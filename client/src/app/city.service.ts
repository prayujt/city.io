import { Injectable } from '@angular/core';
import { Building } from './building'

@Injectable({
  providedIn: 'root'
})
export class CityService {
  private buildings: Building[][] = [];
  constructor() { }

  createCity(size: number = 9) : CityService {
    for (let i = 0; i < size; i++) {
      this.buildings[i] = [];
      for (let j = 0; j < size; j++) {
        if (i == 4 && j == 4)
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
