package open.meteo.domain.model.enums;

import lombok.Getter;

@Getter
public enum ParameterType {

    // Pollutants
    PM10("pm10"),
    PM2_5("pm2_5"),
    CARBON_MONOXIDE("carbon_monoxide"),
    CARBON_DIOXIDE("carbon_dioxide"),
    NITROGEN_DIOXIDE("nitrogen_dioxide"),
    SULPHUR_DIOXIDE("sulphur_dioxide"),
    OZONE("ozone"),

    // Gases & Aerosols
    AMMONIA("ammonia"),
    METHANE("methane"),
    AEROSOL_OPTICAL_DEPTH("aerosol_optical_depth"),
    DUST("dust"),

    // UV Index Variables
    UV_INDEX("uv_index"),
    UV_INDEX_CLEAR_SKY("uv_index_clear_sky");

    private final String name;

    ParameterType(String name) {
        this.name = name;
    }

}
