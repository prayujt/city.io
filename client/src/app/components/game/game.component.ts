import { Component } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';
import { MatSnackBar } from '@angular/material/snack-bar';
import { MatProgressBarModule } from '@angular/material/progress-bar';

import { environment } from '../../../environments/environment';
import { CookieService } from 'ngx-cookie-service';
import { Building } from 'src/app/building';
import { CityService } from '../../city.service';

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
        private cityService: CityService,
        private _snackBar: MatSnackBar,
    ) {
        this.createCity();
    }
    public loggedIn: boolean = false;
    private ID: string = '';

    public ngOnInit(): void {
        this.ID = this.getID();
        if (this.ID != '') {
            this.http
                .get<any>(
                    `http://${environment.API_HOST}:${environment.API_PORT}/sessions/${this.ID}`
                )
                .subscribe((response) => {
                    if (!response.status) {
                        this.router.navigate(['login']);
                    }
                });
            this.http
                .get<any>(
                    `http://${environment.API_HOST}:${environment.API_PORT}/cities/${this.ID}/buildings`
                )
                .subscribe((response) => {
                    this.cityService.setBuildings(response.buildings);
                });
            setInterval(() => {
                this.http
                .get<any>(
                    `http://${environment.API_HOST}:${environment.API_PORT}/cities/${this.ID}/buildings`
                )
                .subscribe((response) => {
                    this.cityService.setBuildings(response.buildings);
                });
            }, 1000);
        } else {
            this.router.navigate(['login']);
        }
    }

    public logOut(): void {
        this.http
            .post<any>(
                `http://${environment.API_HOST}:${environment.API_PORT}/sessions/logout`,
                {
                    sessionId: this.ID,
                }
            )
            .subscribe((response) => {
                if (response.status) {
                    this._snackBar.open('Log out successful!', 'Close', {
                        duration: 2000,
                    });
                    this.loggedIn = false;
                    this.router.navigate(['login']).then(() => {
                        window.location.reload();
                        this.cookieService.delete('cookie');
                    });
                    this.cookieService.delete('cookie')
                } else {
                    this._snackBar.open('Could not log out!', 'Close', {
                        duration: 2000,
                    });
                }
            });
    }

    public getID(): string {
        console.log(this.cookieService.get('cookie'));
        return this.cookieService.get('cookie');
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

    onTileClick(row: number, column: number) {
        // show buildable buildings

        // change CSS for selected tile
        const tiles = Array.from(document.querySelectorAll('.tile'));
        tiles.forEach(tile => tile.classList.remove('selected'));
        const tile = document.querySelector(`td[id="${row} ${column}"]`);
        tile?.classList.add('selected');

        // update sidebar stats
        this.http
            .get<any>(
                `http://${environment.API_HOST}:${environment.API_PORT}/cities/${this.ID}/buildings/${row}/${column}`
            )
            .subscribe((response => {
                this.buildingType = response.buildingType;
                this.buildingLevel = response.buildingLevel;
                this.buildingProduction = response.buildingProduction;
                this.happinessChange = response.happinessChange;
                this.startTime = response.startTime;
                this.endTime = response.endTime;
            }))
        this.clicked = true;
        if (this.startTime != null && this.endTime != null) {
            this.progBar = true;
            this.startUnix = Date.parse(this.startTime);
            this.endUnix = Date.parse(this.endTime);
        }
        else this.progBar = false;
        console.log("clicky")
    }
}
