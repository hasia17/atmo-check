import { Component, OnInit, inject } from '@angular/core';
import { AggregatorService } from '../../services/aggregator.service';
import { AggregatedData } from '../../models/aggregated-data.model';

@Component({
  selector: 'app-map',
  standalone: true,
  template: '<p>Mapa</p>'
})
export class MapComponent implements OnInit {
  private aggregatorService = inject(AggregatorService);
  data: AggregatedData[] = [];

  ngOnInit(): void {
    this.aggregatorService.getAll().subscribe(data => {
      this.data = data;
      console.log('Dane z agregatora:', data);
    });
  }
}
