import { Component, Input, Output, EventEmitter } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { Router } from '@angular/router';
import { MatSnackBar } from '@angular/material/snack-bar';
import { FormControl } from '@angular/forms';
import { CookieService } from 'ngx-cookie-service';

import { Observable, of } from 'rxjs';
import { map, startWith } from 'rxjs/operators';
import { ConstructableService } from '../../../services/constructable.service';
import { Constructable } from 'src/app/services/constructable';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../../../environments/environment';

@Component({
    selector: 'game-sidebar',
    templateUrl: './sidebar.component.html',
    styleUrls: ['./sidebar.component.css'],
})
export class SidebarComponent {
    @Input() row!: number;
    @Input() column!: number;
    @Input() sessionId!: string;
    @Input() isOwner!: boolean;
    @Output() buildBuilding: EventEmitter<string> = new EventEmitter<string>();

    constructor(
        private http: HttpClient,
        private router: Router,
        private cookieService: CookieService,
        private _snackBar: MatSnackBar,
        private dialog: MatDialog
    ) {}

    cityOwner!: string;
    cityName!: string;
    playerBalance: number = 0;
    population!: number;
    populationCapacity!: number;

    buildingType!: string;
    buildingLevel!: number;
    buildingProduction!: number;
    happinessChange!: number;
    startTime!: string;
    endTime!: string;
    clicked: boolean = false;
    progBar: boolean = false;
    constructableBuildings: Constructable[] = [];
    constructableService: ConstructableService = new ConstructableService();

    panelOpenState: boolean = false;
    startUnix!: number;
    endUnix!: number;
    progress!: number;
    remaining!: number;
    dd!: number;
    hh!: number;
    mm!: number;
    ss!: number;

    public ngOnInit(): void {
        // TODO: replace with get request to constructable buildings from database
        this.constructableBuildings = this.constructableService.constructables;

        setInterval(() => {
            if (this.buildingType != '') this.clicked = true;

            let parameter = '';
            let cityName = this.cookieService.get('cityName');
            if (cityName != '') {
                parameter = `?cityName=${encodeURIComponent(cityName)}`;
            }

            this.http
                .get<any>(
                    `${environment.API_HOST}/cities/${this.sessionId}/buildings/${this.row}/${this.column}${parameter}`
                )
                .subscribe((response) => {
                    this.buildingType = response.buildingType;
                    this.buildingLevel = response.buildingLevel;
                    this.buildingProduction = response.buildingProduction;
                    this.happinessChange = response.happinessChange;
                    this.startTime = response.startTime;
                    this.endTime = response.endTime;
                });

            if (this.startTime != '' && this.endTime != '') {
                this.progBar = true;
                this.startUnix = Date.parse(this.startTime);
                this.endUnix = Date.parse(this.endTime);
                let time: number = Date.now();
                this.progress =
                    ((time - this.startUnix) /
                        (this.endUnix - this.startUnix)) *
                    100;
                if (this.progress > 100) this.progBar = false;

                this.remaining = this.endUnix - time;
                this.ss = this.remaining / 1000;
                this.dd = Math.floor(this.ss / 86400);
                this.ss %= 86400;
                this.hh = Math.floor(this.ss / 3600);
                this.ss %= 3600;
                this.mm = Math.floor(this.ss / 60);
                this.ss %= 60;
                this.ss = Math.floor(this.ss);
            } else this.progBar = false;
        }, 100);

        setInterval(() => {
            let parameter = '';
            let cityName = this.cookieService.get('cityName');
            if (cityName != '') {
                parameter = `?cityName=${encodeURIComponent(cityName)}`;
            }

            this.http
                .get<any>(
                    `${environment.API_HOST}/cities/${this.sessionId}${parameter}`
                )
                .subscribe((response) => {
                    this.cityOwner = response.cityOwner;
                    this.cityName = response.cityName;
                    this.playerBalance = response.playerBalance;
                    this.population = response.population;
                    this.populationCapacity = response.populationCapacity;
                });
        }, 500);
    }

    public logOut(): void {
        this.http
            .post<any>(`${environment.API_HOST}/sessions/logout`, {
                sessionId: this.sessionId,
            })
            .subscribe((response) => {
                if (response.status) {
                    this._snackBar.open('Log out successful!', 'Close', {
                        duration: 2000,
                    });
                    this.router.navigate(['login']);
                    this.cookieService.delete('sessionId');
                } else {
                    this._snackBar.open('Could not log out!', 'Close', {
                        duration: 2000,
                    });
                }
            });
    }

    public upgrade(): void {
        let parameter = '';
        let cityName = this.cookieService.get('cityName');
        if (cityName != '') {
            parameter = `?cityName=${encodeURIComponent(cityName)}`;
        }

        this.http
            .post<any>(
                `${environment.API_HOST}/cities/${this.sessionId}/upgradeBuilding${parameter}`,
                {
                    cityRow: this.row,
                    cityColumn: this.column,
                }
            )
            .subscribe((response) => {
                if (response.status) {
                    this._snackBar.open('Upgrade started!', 'Close', {
                        duration: 2000,
                    });
                } else {
                    this._snackBar.open('Error upgrading building!', 'Close', {
                        duration: 2000,
                    });
                }
            });
    }

    public openCityDialog(): void {
        let dialogRef = this.dialog.open(VisitDialogComponent, {
            width: '1000px',
            height: '600px',
        });
    }

    public constructBuilding(buildingType: string): void {
        this.buildBuilding.emit(buildingType);
    }

    public editCityName(): void {
        let dialogRef = this.dialog.open(CityNameChangeDialogComponent, {
            width: '400px',
            height: '200px',
        });

        console.log(this.cityName);
        dialogRef.afterClosed().subscribe((result) => {
            console.log(this.cityName);
            if (result != '' && result != undefined) {
                this.http
                    .post<any>(
                        `${environment.API_HOST}/cities/${this.sessionId}/updateName`,
                        {
                            cityNameOriginal: this.cityName,
                            cityNameNew: result,
                        }
                    )
                    .subscribe((response) => {
                        if (response.status) {
                            this.cookieService.set('cityName', result);
                            this._snackBar.open('City Name Changed!', 'Close', {
                                duration: 2000,
                            });
                        } else {
                            this._snackBar.open(
                                'Error Changing City Name!',
                                'Close',
                                {
                                    duration: 2000,
                                }
                            );
                        }
                    });
            }
        });
    }
}

export interface City {
    cityOwner: string;
    cityName: string;
}

@Component({
    selector: 'visit-dialog',
    templateUrl: './visit-dialog.html',
})
export class VisitDialogComponent {
    cities: City[] = [];
    cityCtrl = new FormControl('');
    filteredCities!: Observable<City[]>;

    towns: City[] = [];
    townCtrl = new FormControl('');
    filteredTowns!: Observable<City[]>;

    constructor(public http: HttpClient, private cookieService: CookieService) {
        this.http
            .get<any>(`${environment.API_HOST}/cities`)
            .subscribe((response) => {
                this.cities = response;
                this.filteredCities = of(this.cities);
                this.filteredCities = this.cityCtrl.valueChanges.pipe(
                    startWith(''),
                    map((city) =>
                        city ? this._filterCities(city) : this.cities.slice()
                    )
                );
            });

        this.http
            .get<any>(`${environment.API_HOST}/towns`)
            .subscribe((response) => {
                this.towns = response;
                this.filteredTowns = of(this.towns);
                this.filteredTowns = this.townCtrl.valueChanges.pipe(
                    startWith(''),
                    map((town) =>
                        town ? this._filterTowns(town) : this.towns.slice()
                    )
                );
            });
    }
    private _filterCities(value: string): any[] {
        const filterValue = value.toLowerCase();

        return this.cities.filter(
            (city) =>
                city.cityOwner.toLowerCase().includes(filterValue) ||
                city.cityName.toLowerCase().includes(filterValue)
        );
    }

    private _filterTowns(value: string): any[] {
        const filterValue = value.toLowerCase();

        return this.towns.filter(
            (town) =>
                town.cityOwner.toLowerCase().includes(filterValue) ||
                town.cityName.toLowerCase().includes(filterValue)
        );
    }

    public visitCity(cityName: string) {
        this.cookieService.set('cityName', cityName);
    }
}

@Component({
    selector: 'city-name-change-dialog',
    templateUrl: './city-name-change-dialog.html',
})
export class CityNameChangeDialogComponent {
    constructor() {}
}
