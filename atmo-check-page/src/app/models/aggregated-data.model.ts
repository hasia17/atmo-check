export type ParamType = 'PM10' | 'PM2_5' | 'CO' | 'CO2' | 'NO2' | 'SO2' | 'O3' | 'CH4';

export interface Parameter {
  id: number;
  description: string;
  unit: string;
  value: number;
  type: ParamType;
}

export interface AggregatedData {
  voivodeship: string;
  parameters: Parameter[];
  timestamp: string;
}
