import { Injectable } from '@angular/core';
import { Constructable } from './constructable';
import { HttpClient } from '@angular/common/http';
import { environment } from 'src/environments/environment';
@Injectable({
  providedIn: 'root'
})
export class ConstructableService {
  constructables: Constructable[] = [];
  public sessionId: string = '';

  constructor(
    private http: HttpClient
  ) {
    http
      .get<any>(
        `${environment.API_HOST}/sessions/cities/buildings`
      ).subscribe((response) => {
        this.constructables = this.setConstructables(response.buildings);
      });
  }

  setConstructables(buildings: Map<string, number[]>): Constructable[] {
    let i = 0;
    let constructables: Constructable[] = [];
    buildings.forEach(function(value, key) {
      constructables[i] = new Constructable(key, value);
      i++;
    });
    return constructables;
  }
}
