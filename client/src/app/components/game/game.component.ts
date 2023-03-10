import { Component } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';
import { MatSnackBar } from '@angular/material/snack-bar';

import { environment } from '../../../environments/environment';
import { CookieService } from 'ngx-cookie-service';
import { Building } from '../../services/building';
import { CityService } from '../../services/city.service';
import { SidebarComponent } from './sidebar/sidebar.component';

@Component({
    selector: 'app-game',
    templateUrl: './game.component.html',
    styleUrls: ['./game.component.css'],
})
export class GameComponent {
    constructor(
        private http: HttpClient,
        private router: Router,
        private cookieService: CookieService,
        private _snackBar: MatSnackBar,
        private cityService: CityService
    ) {
        this.createCity();
    }
    public row: number = 0;
    public column: number = 0;
    public sessionId: string = '';
    public isOwner!: boolean;
    public homeCity!: string;

    public ngOnInit(): void {
        // this.cookieService.delete('cityName');
        this.sessionId = this.getID();
        if (this.sessionId != '') {
            this.http
                .get<any>(
                    `${environment.API_HOST}/sessions/${this.sessionId}`
                )
                .subscribe((response) => {
                    if (!response.status) {
                        this.router.navigate(['login']);
                    }
                });
            let parameter = '';
            if (this.cookieService.get('cityName') != '') {
                parameter = `?cityName=${encodeURIComponent(
                    this.cookieService.get('cityName')
                )}`;
            }

            this.http
                .get<any>(
                    `${environment.API_HOST}/cities/${this.sessionId}/buildings${parameter}`
                )
                .subscribe((response) => {
                    this.cityService.setBuildings(response.buildings);
                });

            setInterval(() => {
                let parameter = '';
                if (this.cookieService.get('cityName') != '') {
                    parameter = `?cityName=${encodeURIComponent(
                        this.cookieService.get('cityName')
                    )}`;
                }
                this.http
                    .get<any>(
                        `${environment.API_HOST}/cities/${this.sessionId}/buildings${parameter}`
                    )
                    .subscribe(async (response) => {
                        this.cityService.setBuildings(response.buildings);
                        this.isOwner = await response.isOwner;
                        if (this.isOwner && this.cookieService.get('cityName') != '') {
                            this.homeCity = this.cookieService.get('cityName');
                            console.log('Home: ' + this.homeCity);
                        }
                    });
            }, 250);
        } else {
            this.router.navigate(['login']);
        }
    }

    public getID(): string {
        return this.cookieService.get('sessionId');
    }

    createCity(): GameComponent {
        this.cityService.createCity();
        return this;
    }

    get buildings(): Building[][] {
        return this.cityService.getBuildings();
    }

    buildingType!: string;
    buildingLevel!: number;
    buildingProduction!: number;
    happinessChange!: number;
    startTime!: string;
    endTime!: string;
    clicked: boolean = false;
    progBar: boolean = false;

    startUnix!: number;
    endUnix!: number;
    progress!: number;
    interval: any;
    remaining!: number;
    hh!: number;
    mm!: number;
    ss!: number;

    onTileClick(row: number, column: number) {
        // change CSS for selected tile
        const tiles = Array.from(document.querySelectorAll('.tile'));
        tiles.forEach((tile) => tile.classList.remove('selected'));
        const tile = document.querySelector(`td[id="${row} ${column}"]`);
        tile?.classList.add('selected');

        this.row = row;
        this.column = column;
    }

    onBuildingClick(buildingType: string): void {
        this.http
            .post<any>(
                `${environment.API_HOST}/cities/${this.sessionId}/createBuilding`,
                {
                    buildingType: buildingType,
                    cityRow: this.row,
                    cityColumn: this.column,
                }
            )
            .subscribe((response) => {
                if (!response.status) {
                    this._snackBar.open(
                        'Building could not be constructed!',
                        'Close',
                        {
                            duration: 2000,
                        }
                    );
                } else {
                    this._snackBar.open('Building constructed!', 'Close', {
                        duration: 2000,
                    });
                }
            });
    }
}
