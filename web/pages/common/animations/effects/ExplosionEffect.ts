import Phaser from 'phaser';
import { AnimationConfig, AnimationVisualConfig } from '../AnimationConfig';

/**
 * Explosion particle effect.
 * Creates a burst of particles at a target location with intensity scaling.
 */
export class ExplosionEffect {
  private scene: Phaser.Scene;
  private x: number;
  private y: number;
  private intensity: number;

  constructor(scene: Phaser.Scene, x: number, y: number, intensity: number = 1) {
    this.scene = scene;
    this.x = x;
    this.y = y;
    this.intensity = Math.max(0, intensity);
  }

  /**
   * Play the explosion effect.
   * Returns a promise that resolves when animation completes.
   */
  public play(): Promise<void> {
    return new Promise((resolve) => {
      if (AnimationConfig.EXPLOSION_DURATION === 0) {
        // Instant mode
        resolve();
        return;
      }

      // Calculate particle count based on intensity
      const baseCount = AnimationVisualConfig.EXPLOSION_PARTICLE_COUNT;
      const intensityBonus = this.intensity * AnimationVisualConfig.EXPLOSION_PARTICLES_PER_DAMAGE;
      const particleCount = Math.min(
        baseCount + intensityBonus,
        AnimationVisualConfig.EXPLOSION_PARTICLE_MAX
      );

      // Create particle emitter configuration
      const particles = this.scene.add.particles(this.x, this.y, 'particle', {
        speed: { min: 50, max: 150 },
        angle: { min: 0, max: 360 },
        scale: { start: 1.0, end: 0 },
        alpha: { start: 1, end: 0 },
        lifespan: AnimationConfig.EXPLOSION_DURATION,
        quantity: particleCount,
        blendMode: Phaser.BlendModes.ADD,
        tint: [
          AnimationVisualConfig.EXPLOSION_COLOR_MIN,
          AnimationVisualConfig.EXPLOSION_COLOR_MAX,
        ],
        emitting: false,
      });

      particles.setDepth(20);

      // Emit burst
      particles.explode(particleCount);

      // Clean up after animation
      this.scene.time.delayedCall(AnimationConfig.EXPLOSION_DURATION, () => {
        particles.destroy();
        resolve();
      });
    });
  }

  /**
   * Create multiple simultaneous explosions (for splash damage).
   * Returns a promise that resolves when all explosions complete.
   */
  public static playMultiple(
    scene: Phaser.Scene,
    positions: { x: number; y: number; intensity: number }[]
  ): Promise<void> {
    const explosions = positions.map(
      (pos) => new ExplosionEffect(scene, pos.x, pos.y, pos.intensity).play()
    );
    return Promise.all(explosions).then(() => {});
  }
}
