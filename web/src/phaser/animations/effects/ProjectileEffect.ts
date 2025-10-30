import Phaser from 'phaser';
import { AnimationConfig, AnimationVisualConfig } from '../AnimationConfig';

/**
 * Projectile effect that flies from source to target with an arc trajectory.
 * Returns a promise that resolves when the projectile reaches its target.
 */
export class ProjectileEffect extends Phaser.GameObjects.Graphics {
  private arcHeight: number;
  private startX: number;
  private startY: number;
  private endX: number;
  private endY: number;

  constructor(
    scene: Phaser.Scene,
    fromX: number,
    fromY: number,
    toX: number,
    toY: number
  ) {
    super(scene);

    this.startX = fromX;
    this.startY = fromY;
    this.endX = toX;
    this.endY = toY;

    // Calculate arc height based on distance
    const distance = Phaser.Math.Distance.Between(fromX, fromY, toX, toY);
    this.arcHeight = distance * AnimationVisualConfig.PROJECTILE_ARC_HEIGHT;

    // Set depth to render above units
    this.setDepth(15);

    // Draw the projectile as a small circle
    this.fillStyle(0x333333, 1.0);
    this.fillCircle(0, 0, 4);

    // Position at start
    this.setPosition(fromX, fromY);

    scene.add.existing(this);
  }

  /**
   * Play the projectile animation from start to end.
   * Returns a promise that resolves when animation completes.
   */
  public play(): Promise<void> {
    return new Promise((resolve) => {
      if (AnimationConfig.PROJECTILE_DURATION === 0) {
        // Instant mode
        this.destroy();
        resolve();
        return;
      }

      // Create arc trajectory using a parabolic path
      this.scene.tweens.add({
        targets: this,
        x: this.endX,
        duration: AnimationConfig.PROJECTILE_DURATION,
        ease: 'Linear',
        onUpdate: (tween) => {
          // Calculate arc based on progress
          const progress = tween.progress;
          const arcOffset = Math.sin(progress * Math.PI) * this.arcHeight;
          const currentY = Phaser.Math.Linear(this.startY, this.endY, progress);
          this.y = currentY - arcOffset;
        },
        onComplete: () => {
          this.destroy();
          resolve();
        },
      });
    });
  }
}
