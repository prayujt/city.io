import { Injectable } from '@angular/core';
import { Building } from './building';

@Injectable({
    providedIn: 'root',
})
export class CityService {
    private buildings: Building[][] = [[]];
    private availableBuildings: Building[] = [
        new Building(-1, )
    ];
    troops: Map<number, number> = new Map([
        [1, 1000],
        [2, 5000],
        [3, 10000],
        [4, 25000],
        [5, 50000],
        [6, 100000],
        [7, 250000],
        [8, 500000],
        [9, 1000000],
        [10, 10000000]
    ]);
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

    getAvailableBuildings(): Building[] {
        return this.availableBuildings;
    }

    updateIcons() {
        for (let i = 0; i < 9; i++) {
            for (let j = 0; j < 13; j++) {
                this.buildings[i][j].setIcon();
            }
        }
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
        let maxTrainCount: number = 0
        if (buildings == null) {
            for (let i = 0; i < 9; i++) {
                for (let j = 0; j < 13; j++) {
                    this.buildings[i][j].level = 0;
                    this.buildings[i][j].type = '';
                    this.buildings[i][j].setIcon();
                }
            }
        } else {
            for (let i = 0; i < 9; i++) {
                for (let j = 0; j < 13; j++) {
                    let matched = false;
                    for (let k = 0; k < buildings.length; k++) {
                        if (
                            buildings[k].cityRow == i &&
                            buildings[k].cityColumn == j
                        ) {
                            this.buildings[i][j].level =
                                buildings[k].buildingLevel;
                            this.buildings[i][j].type =
                                buildings[k].buildingType;
                            if (this.buildings[i][j].type == 'Barracks') {
                                maxTrainCount += this.troops.get(buildings[k].buildingLevel) as number;
                            }
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
        return maxTrainCount;
    }
}
