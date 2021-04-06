import { LiveChannelSupport, LiveChannelConfig } from '@grafana/data';

interface MeasurementChannel {
  config: LiveChannelConfig;
}

// HACK -- this is copied from:
// from 'app/features/live/measurements/measurementsSupport';
// likely should find the right pattern and export that
export class LiveMeasurementsSupport implements LiveChannelSupport {
  private cache: Record<string, MeasurementChannel> = {};

  /**
   * Get the channel handler for the path, or throw an error if invalid
   */
  getChannelConfig(path: string): LiveChannelConfig | undefined {
    let c = this.cache[path];
    if (!c) {
      c = this.cache[path] = {
        config: {
          path,
          canPublish: () => true,
        },
      };
    }
    return c.config;
  }

  /**
   * Return a list of supported channels
   */
  getSupportedPaths(): LiveChannelConfig[] {
    // this should ask the server what channels it has seen
    return [];
  }
}
