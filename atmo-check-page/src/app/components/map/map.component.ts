import { Component, OnInit, inject, signal, computed } from '@angular/core';
import { MatCardModule } from '@angular/material/card';
import { MatButtonToggleModule } from '@angular/material/button-toggle';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { AggregatorService } from '../../services/aggregator.service';
import { AggregatedData, Parameter, ParamType } from '../../models/aggregated-data.model';
import { VOIVODESHIPS } from './poland-map';

@Component({
  selector: 'app-map',
  standalone: true,
  imports: [MatCardModule, MatButtonToggleModule, MatIconModule, MatProgressSpinnerModule],
  templateUrl: './map.component.html',
  styleUrl: './map.component.scss'
})
export class MapComponent implements OnInit {
  private aggregatorService = inject(AggregatorService);

  data = signal<AggregatedData[]>([]);
  loading = signal(true);
  selectedParam = signal<ParamType>('PM10');

  readonly params: ParamType[] = ['PM10', 'PM2_5', 'CO', 'CO2', 'NO2', 'SO2', 'O3', 'CH4'];
  readonly voivodeships = VOIVODESHIPS;

  private values = computed(() => {
    const map = new Map<string, number>();
    for (const item of this.data()) {
      const param = item.parameters.find(p => p.type === this.selectedParam());
      if (param) map.set(item.voivodeship, param.value);
    }
    return map;
  });

  legendRange = computed(() => {
    const vals = Array.from(this.values().values());
    if (!vals.length) return { min: 0, max: 0 };
    return { min: Math.min(...vals), max: Math.max(...vals) };
  });

  selectedParamInfo = computed<Parameter | null>(() => {
    const first = this.data()[0];
    if (!first) return null;
    return first.parameters.find(p => p.type === this.selectedParam()) ?? null;
  });

  ngOnInit(): void {
    this.aggregatorService.getAll().subscribe(data => {
      this.data.set(data);
      this.loading.set(false);
    });
  }

  getColor(voivodeshipId: string): string {
    const vals = this.values();
    const value = vals.get(voivodeshipId);
    if (value === undefined) return '#cccccc';

    const { min, max } = this.legendRange();
    const ratio = max === min ? 0.5 : (value - min) / (max - min);

    if (ratio < 0.5) {
      return interpolate('#4caf50', '#ffeb3b', ratio * 2);
    } else {
      return interpolate('#ffeb3b', '#f44336', (ratio - 0.5) * 2);
    }
  }

  getValue(voivodeshipId: string): string {
    const value = this.values().get(voivodeshipId);
    if (value === undefined) return '';
    return value.toFixed(1);
  }

  onParamChange(param: ParamType): void {
    this.selectedParam.set(param);
  }
}

function interpolate(from: string, to: string, t: number): string {
  const f = hexToRgb(from);
  const t2 = hexToRgb(to);
  const r = Math.round(f.r + (t2.r - f.r) * t);
  const g = Math.round(f.g + (t2.g - f.g) * t);
  const b = Math.round(f.b + (t2.b - f.b) * t);
  return `rgb(${r},${g},${b})`;
}

function hexToRgb(hex: string) {
  return {
    r: parseInt(hex.slice(1, 3), 16),
    g: parseInt(hex.slice(3, 5), 16),
    b: parseInt(hex.slice(5, 7), 16),
  };
}
