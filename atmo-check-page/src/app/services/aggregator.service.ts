import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { AggregatedData } from '../models/aggregated-data.model';

@Injectable({ providedIn: 'root' })
export class AggregatorService {
  private readonly apiUrl = 'http://localhost:8082/aggregatedData';

  constructor(private http: HttpClient) {}

  getAll(): Observable<AggregatedData[]> {
    return this.http.get<AggregatedData[]>(this.apiUrl);
  }
}
