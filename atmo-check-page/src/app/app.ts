import { Component } from '@angular/core';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MapComponent } from './components/map/map.component';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [MatToolbarModule, MapComponent],
  templateUrl: './app.html',
  styleUrl: './app.scss'
})
export class App {}
