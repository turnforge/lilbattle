import Phaser from 'phaser';
import { AnimationConfig } from '../AnimationConfig';

/**
 * Capture/occupation effect.
 * Shows a visual indication that a tile or building has been captured.
 * Currently implements a simple color pulse effect.
 */
export class CaptureEffect {
  private scene: Phaser.Scene;
  private x: number;
  private y: number;

  constructor(scene: Phaser.Scene, x: number, y: number) {
    this.scene = scene;
    this.x = x;
    this.y = y;
  }

  /**
   * Play the capture effect.
   * Returns a promise that resolves when animation completes.
   */
  public play(): Promise<void> {
    return new Promise((resolve) => {
      if (AnimationConfig.CAPTURE_DURATION === 0) {
        // Instant mode
        resolve();
        return;
      }

      // Create a circular pulse effect
      const circle = this.scene.add.circle(this.x, this.y, 10, 0xffaa00, 0.8);
      circle.setDepth(19);

      // Animate: expand and fade out
      this.scene.tweens.add({
        targets: circle,
        radius: 40,
        alpha: 0,
        duration: AnimationConfig.CAPTURE_DURATION,
        ease: 'Cubic.easeOut',
        onComplete: () => {
          circle.destroy();
          resolve();
        },
      });
    });
  }
}
