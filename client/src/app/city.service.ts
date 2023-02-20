import { Injectable } from '@angular/core';
import { Building } from './building'

@Injectable({
  providedIn: 'root'
})
export class CityService {
  private buildings: Building[][] = [];
  constructor() { }

  createCity() : CityService {
    for (let i = 0; i < 9; i++) {
      this.buildings[i] = [];
      for (let j = 0; j < 13; j++) {
        if (i == 4 && j == 6)
          this.buildings[i][j] = new Building("city_hall");
        else
        this.buildings[i][j] = new Building("");
      }
    }

    return this;
  }

  getBuildings() : Building[][] {
    return this.buildings;
  }
}
