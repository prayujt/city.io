import { Component } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
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
    public jwtToken: string = '';
    public isOwner: boolean = true;

    public ngOnInit(): void {
        this.cookieService.delete('cityName');
        this.jwtToken = this.getID();
        if (this.jwtToken != '') {
            let headers = new HttpHeaders();
            headers = headers.append('Token', this.jwtToken);
            this.http
                .get<any>(`${environment.API_HOST}/sessions/validate`, {
                    headers,
                })
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
                    `${environment.API_HOST}/cities/buildings${parameter}`,
                    { headers }
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
                        `${environment.API_HOST}/cities/buildings${parameter}`,
                        { headers }
                    )
                    .subscribe((response) => {
                        this.cityService.setBuildings(response.buildings);
                        this.isOwner = response.isOwner;
                    });
            }, 250);
        } else {
            this.router.navigate(['login']);
        }
    }

    public getID(): string {
        return this.cookieService.get('jwtToken');
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
        const tiles = Array.from(document.querySelectorAll('.tile'));
        tiles.forEach((tile) => tile.classList.remove('selected'));
        const tile = document.querySelector(`td[id="${row} ${column}"]`);
        tile?.classList.add('selected');

        this.row = row;
        this.column = column;
    }

    onBuildingClick(buildingType: string): void {
        let parameter = '';
        let cityName = this.cookieService.get('cityName');
        if (cityName != '') {
            parameter = `?cityName=${encodeURIComponent(cityName)}`;
        }

        let headers = new HttpHeaders();
        headers = headers.append('Token', this.jwtToken);

        this.http
            .post<any>(
                `${environment.API_HOST}/cities/createBuilding${parameter}`,
                {
                    buildingType: buildingType,
                    cityRow: this.row,
                    cityColumn: this.column,
                },
                { headers }
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
