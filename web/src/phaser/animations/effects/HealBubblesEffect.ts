import Phaser from 'phaser';
import { AnimationConfig, AnimationVisualConfig } from '../AnimationConfig';

/**
 * Healing bubble particle effect.
 * Creates rising bubbles at a target location to indicate healing.
 */
export class HealBubblesEffect {
  private scene: Phaser.Scene;
  private x: number;
  private y: number;
  private amount: number;

  constructor(scene: Phaser.Scene, x: number, y: number, amount: number = 1) {
    this.scene = scene;
    this.x = x;
    this.y = y;
    this.amount = Math.max(0, amount);
  }

  /**
   * Play the healing bubble effect.
   * Returns a promise that resolves when animation completes.
   */
  public play(): Promise<void> {
    return new Promise((resolve) => {
      if (AnimationConfig.HEAL_DURATION === 0) {
        // Instant mode
        resolve();
        return;
      }

      // Create particle emitter for rising bubbles
      const particles = this.scene.add.particles(this.x, this.y, 'particle', {
        speed: { min: 20, max: 40 },
        angle: { min: -100, max: -80 }, // Upward direction
        scale: { start: 0.5, end: 0.2 },
        alpha: { start: 0.8, end: 0 },
        lifespan: AnimationConfig.HEAL_DURATION,
        quantity: AnimationVisualConfig.HEAL_PARTICLE_COUNT,
        gravityY: -50, // Float upward
        blendMode: Phaser.BlendModes.ADD,
        tint: [
          AnimationVisualConfig.HEAL_COLOR_MIN,
          AnimationVisualConfig.HEAL_COLOR_MAX,
        ],
        emitting: false,
      });

      particles.setDepth(20);

      // Emit burst
      particles.explode(AnimationVisualConfig.HEAL_PARTICLE_COUNT);

      // Clean up after animation
      this.scene.time.delayedCall(AnimationConfig.HEAL_DURATION, () => {
        particles.destroy();
        resolve();
      });
    });
  }
}
