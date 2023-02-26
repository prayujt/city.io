import { Injectable } from '@angular/core';
import { Building } from './building';

@Injectable({
    providedIn: 'root',
})
export class CityService {
    private buildings: Building[][] = [[]];
    constructor() {}

    createCity(): CityService {
        for (let i = 0; i < 9; i++) {
            this.buildings[i] = [];
            for (let j = 0; j < 13; j++) {
                this.buildings[i][j] = new Building(0, '', i, j);
            }
        }

        return this;
    }

    getBuildings(): Building[][] {
        return this.buildings;
    }

    setBuildings(
        buildings: Array<{
            buildingLevel: number;
            buildingName: string;
            buildingType: string;
            cityColumn: number;
            cityRow: number;
        }>
    ) {
        if (buildings == null) {
            for (let i = 0; i < 9; i++) {
                for (let j = 0; j < 13; j++) {
                    this.buildings[i][j].level = 0;
                    this.buildings[i][j].type = '';
                }
            }
        }
        for (let i = 0; i < 9; i++) {
            for (let j = 0; j < 13; j++) {
                let matched = false;
                for (let k = 0; k < buildings.length; k++) {
                    if (
                        buildings[k].cityRow == i &&
                        buildings[k].cityColumn == j
                    ) {
                        this.buildings[i][j].level = buildings[k].buildingLevel;
                        this.buildings[i][j].type = buildings[k].buildingType;
                        matched = true;
                    }
                    if (!matched) {
                        this.buildings[i][j].level = 0;
                        this.buildings[i][j].type = '';
                    }
                    this.buildings[i][j].setIcon();
                }
            }
        }
    }
}
