package batching;

import java.time.LocalDateTime;

public class Medicion {
    public Medicion (LocalDateTime datetime, String sensor, String sector, int presion) {
        this.datetime = datetime;
        this.sensor = sensor;
        this.sector = sector;
        this.presion = presion;
    }

    private LocalDateTime datetime;
    public LocalDateTime getDatetime() {
        return datetime;
    }

    private String sensor; 
    public String getSensor() {
        return sensor;
    }

    private String sector;
    public String getSector() {
        return sector;
    }

    private int presion;
    public int getPresion() {
        return presion;
    }
}
