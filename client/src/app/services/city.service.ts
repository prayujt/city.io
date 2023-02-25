import { Injectable } from '@angular/core';
import { Building } from './building';

@Injectable({
    providedIn: 'root',
})
export class CityService {
    private buildings: Building[][] = [];
    constructor() {}

    createCity(): CityService {
        for (let i = 0; i < 9; i++) {
            this.buildings[i] = [];
            for (let j = 0; j < 13; j++) {
                this.buildings[i][j] = new Building(0, '', '', i, j);
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
        this.createCity();
        for (let i = 0; i < buildings.length; i++) {
            this.buildings[buildings[i].cityRow][
                buildings[i].cityColumn
            ].level = buildings[i].buildingLevel;
            this.buildings[buildings[i].cityRow][buildings[i].cityColumn].name =
                buildings[i].buildingName;
            this.buildings[buildings[i].cityRow][buildings[i].cityColumn].type =
                buildings[i].buildingType;
            this.buildings[buildings[i].cityRow][
                buildings[i].cityColumn
            ].setIcon();
        }
    }
}
