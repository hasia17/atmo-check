package aggregator.service;

import aggregator.model.Voivodeship;

import java.util.HashMap;
import java.util.Map;

/**
 * Utility class for determining voivodeship boundaries based on coordinates
 */
public class VoivodeshipBounds {

    public static class GeographicBounds {
        private final double minLat;
        private final double maxLat;
        private final double minLon;
        private final double maxLon;

        public GeographicBounds(double minLat, double maxLat, double minLon, double maxLon) {
            this.minLat = minLat;
            this.maxLat = maxLat;
            this.minLon = minLon;
            this.maxLon = maxLon;
        }

        public boolean contains(double lat, double lon) {
            return lat >= minLat && lat <= maxLat && lon >= minLon && lon <= maxLon;
        }
    }

    private static final Map<Voivodeship, GeographicBounds> BOUNDS = new HashMap<>();

    static {
        BOUNDS.put(Voivodeship.DOLNOSLASKIE, new GeographicBounds(50.0, 51.8, 15.0, 17.8));
        BOUNDS.put(Voivodeship.KUJAWSKO_POMORSKIE, new GeographicBounds(52.0, 53.8, 17.0, 19.8));
        BOUNDS.put(Voivodeship.LUBELSKIE, new GeographicBounds(50.2, 51.6, 21.8, 24.1));
        BOUNDS.put(Voivodeship.LUBUSKIE, new GeographicBounds(51.2, 52.9, 14.1, 16.2));
        BOUNDS.put(Voivodeship.LODZKIE, new GeographicBounds(51.0, 52.5, 18.2, 20.6));
        BOUNDS.put(Voivodeship.MALOPOLSKIE, new GeographicBounds(49.1, 50.8, 18.8, 21.3));
        BOUNDS.put(Voivodeship.MAZOWIECKIE, new GeographicBounds(51.4, 53.4, 19.1, 22.9));
        BOUNDS.put(Voivodeship.OPOLSKIE, new GeographicBounds(50.0, 51.1, 17.0, 18.9));
        BOUNDS.put(Voivodeship.PODKARPACKIE, new GeographicBounds(49.0, 50.9, 21.0, 23.0));
        BOUNDS.put(Voivodeship.PODLASKIE, new GeographicBounds(52.5, 54.4, 22.1, 24.2));
        BOUNDS.put(Voivodeship.POMORSKIE, new GeographicBounds(53.4, 54.8, 16.8, 19.3));
        BOUNDS.put(Voivodeship.SLASKIE, new GeographicBounds(49.8, 50.8, 18.4, 19.8));
        BOUNDS.put(Voivodeship.SWIETOKRZYSKIE, new GeographicBounds(50.1, 51.1, 19.7, 21.5));
        BOUNDS.put(Voivodeship.WARMINSKO_MAZURSKIE, new GeographicBounds(53.2, 54.4, 19.3, 22.6));
        BOUNDS.put(Voivodeship.WIELKOPOLSKIE, new GeographicBounds(51.2, 53.1, 15.6, 18.9));
        BOUNDS.put(Voivodeship.ZACHODNIOPOMORSKIE, new GeographicBounds(52.7, 54.9, 14.1, 16.8));
    }

    /**
     * Determines which voivodeship contains the given coordinates
     *
     * @param lat latitude
     * @param lon longitude
     * @return Voivodeship if found, null otherwise
     */
    public static Voivodeship getVoivodeshipForCoordinates(double lat, double lon) {
        for (Map.Entry<Voivodeship, GeographicBounds> entry : BOUNDS.entrySet()) {
            if (entry.getValue().contains(lat, lon)) {
                return entry.getKey();
            }
        }
        return null;
    }

    /**
     * Checks if given coordinates are within the specified voivodeship
     *
     * @param lat         latitude
     * @param lon         longitude
     * @param voivodeship voivodeship to check
     * @return true if coordinates are within the voivodeship
     */
    public static boolean isInVoivodeship(double lat, double lon, Voivodeship voivodeship) {
        GeographicBounds bounds = BOUNDS.get(voivodeship);
        return bounds != null && bounds.contains(lat, lon);
    }

}
