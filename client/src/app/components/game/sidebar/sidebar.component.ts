import { Component, Input, Inject, Output, EventEmitter } from '@angular/core';

import {
    MatDialog,
    MatDialogRef,
    MAT_DIALOG_DATA,
} from '@angular/material/dialog';
import { Router } from '@angular/router';
import { MatSnackBar } from '@angular/material/snack-bar';
import { FormControl } from '@angular/forms';
import { CookieService } from 'ngx-cookie-service';

import { Observable, of } from 'rxjs';
import { map, startWith } from 'rxjs/operators';
import { ConstructableService } from '../../../services/constructable.service';
import { Constructable } from 'src/app/services/constructable';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { environment } from '../../../../environments/environment';

@Component({
    selector: 'game-sidebar',
    templateUrl: './sidebar.component.html',
    styleUrls: ['./sidebar.component.css'],
})
export class SidebarComponent {
    @Input() row!: number;
    @Input() column!: number;
    @Input() jwtToken!: string;
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
    armySize: number = 0;

    buildingType!: string;
    buildingLevel!: number;
    buildingProduction!: number;
    happinessChange!: number;
    startTime!: string;
    endTime!: string;
    clicked: boolean = false;
    progBar: boolean = false;
    constructableService: ConstructableService = new ConstructableService(
        this.http,
        this.cookieService
    );
    constructableBuildings: Constructable[] = [];

    panelOpenState: boolean = false;
    startUnix!: number;
    endUnix!: number;
    progress!: number;
    remaining!: number;
    dd!: number;
    hh!: number;
    mm!: number;
    ss!: number;

    interval1!: ReturnType<typeof setInterval>;
    interval2!: ReturnType<typeof setInterval>;

    @Input() maxTrainCount: number = 0;

    public ngOnInit(): void {
        this.interval1 = setInterval(() => {
            this.constructableBuildings =
                this.constructableService.getConstructables();
            if (this.buildingType != '') this.clicked = true;

            let parameter = '';
            let cityName = this.cookieService.get('cityName');
            if (cityName != '') {
                parameter = `?cityName=${encodeURIComponent(cityName)}`;
            }

            let headers = new HttpHeaders();
            headers = headers.append('Token', this.jwtToken);
            this.http
                .get<any>(
                    `${environment.API_HOST}/cities/buildings/${this.row}/${this.column}${parameter}`,
                    { headers }
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
        }, 200);

        this.interval2 = setInterval(() => {
            let parameter = '';
            let cityName = this.cookieService.get('cityName');
            if (cityName != '') {
                parameter = `?cityName=${encodeURIComponent(cityName)}`;
            }

            let headers = new HttpHeaders();
            headers = headers.append('Token', this.jwtToken);
            this.http
                .get<any>(`${environment.API_HOST}/cities/stats${parameter}`, {
                    headers,
                })
                .subscribe((response) => {
                    this.cityOwner = response.cityOwner;
                    this.cityName = response.cityName;
                    this.playerBalance = response.playerBalance;
                    this.population = response.population;
                    this.populationCapacity = response.populationCapacity;
                    this.armySize = response.armySize;
                });
        }, 500);
    }

    public ngOnDestroy(): void {
        clearInterval(this.interval1);
        clearInterval(this.interval2);
    }

    public logOut(): void {
        this._snackBar.open('Log out successful!', 'Close', {
            duration: 2000,
        });
        this.router.navigate(['login']);
        this.cookieService.delete('jwtToken');
    }

    public upgrade(): void {
        let parameter = '';
        let cityName = this.cookieService.get('cityName');
        if (cityName != '') {
            parameter = `?cityName=${encodeURIComponent(cityName)}`;
        }

        let headers = new HttpHeaders();
        headers = headers.append('Token', this.jwtToken);
        this.http
            .post<any>(
                `${environment.API_HOST}/cities/upgradeBuilding${parameter}`,
                {
                    cityRow: this.row,
                    cityColumn: this.column,
                },
                { headers }
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
        this.dialog.open(VisitDialogComponent, {
            width: '1000px',
            height: '600px',
        });
    }

    public openAttackDialog(): void {
        this.dialog.open(AttackDialogComponent, {
            width: '1000px',
            height: '600px',
            data: { cityName: this.cityName },
        });
    }

    public openMarchesDialog(): void {
        this.dialog.open(MarchesDialogComponent, {
            width: '1000px',
            height: '600px',
        });
    }

    public openTrainDialog(): void {
        this.dialog.open(TrainDialogComponent, {
            width: '1000px',
            height: '600px',
            data: {
                cityName: this.cityName,
                maxTrainCount: this.maxTrainCount,
            },
        });
    }

    public openBattleLogsDialog(): void {
        this.dialog.open(BattleLogsDialogComponent, {
            width: '1000px',
            height: '600px',
            data: { cityName: this.cityName, cityOwner: this.cityOwner },
        });
    }

    public openTerritoryDialog(): void {
        this.dialog.open(TerritoryDialogComponent, {
            width: '1000px',
            height: '600px',
        });
    }

    public constructBuilding(buildingType: string): void {
        this.buildBuilding.emit(buildingType);
    }

    public openDeleteDialog(): void {
        this.dialog.open(DeleteDialogComponent, {
            width: '210px',
            height: '150px',
            data: { cityRow: this.row, cityColumn: this.column },
        });
    }

    public editCityName(): void {
        let dialogRef = this.dialog.open(CityNameChangeDialogComponent, {
            width: '400px',
            height: '200px',
        });

        dialogRef.afterClosed().subscribe((result) => {
            if (result != '' && result != undefined) {
                let headers = new HttpHeaders();
                headers = headers.append('Token', this.jwtToken);
                this.http
                    .post<any>(
                        `${environment.API_HOST}/cities/updateName`,
                        {
                            cityNameOriginal:
                                this.cookieService.get('cityName'),
                            cityNameNew: result,
                        },
                        { headers }
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

    public goHome(): void {
        this.cookieService.delete('cityName');
    }

    public playAudio(filename: string) {
        let audio = new Audio();
        audio.src = '../../../../assets/audio/' + filename;
        audio.load();
        audio.play();
    }
}

export interface City {
    cityOwner: string;
    cityName: string;
}

export interface CityArmy {
    cityName: string;
    armySize: number;
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
    selector: 'attack-dialog',
    templateUrl: './attack-dialog.html',
    styleUrls: ['./sidebar.component.css'],
})
export class AttackDialogComponent {
    ownedTerritory: CityArmy[] = [];
    territoryCtrl = new FormControl('');
    filteredTerritory!: Observable<CityArmy[]>;

    selectedTerritory!: CityArmy;
    selected: boolean = false;

    armySize: number = 1;

    constructor(
        public dialogRef: MatDialogRef<AttackDialogComponent>,
        public http: HttpClient,
        private cookieService: CookieService,
        @Inject(MAT_DIALOG_DATA) public data: { cityName: string }
    ) {
        let headers = new HttpHeaders();
        headers = headers.append('Token', this.cookieService.get('jwtToken'));
        this.http
            .get<any>(`${environment.API_HOST}/cities/territory`, { headers })
            .subscribe((response) => {
                this.ownedTerritory = response;

                this.filteredTerritory = of(this.ownedTerritory);
                this.filteredTerritory = this.territoryCtrl.valueChanges.pipe(
                    startWith(''),
                    map((territory) =>
                        territory
                            ? this._filterCities(territory)
                            : this.ownedTerritory.slice()
                    )
                );
            });
    }

    public saveCityInfo(selectedTerritory: CityArmy) {
        this.selected = true;
        this.selectedTerritory = selectedTerritory;
    }

    public march() {
        let headers = new HttpHeaders();
        headers = headers.append('Token', this.cookieService.get('jwtToken'));
        this.http
            .post<any>(
                `${environment.API_HOST}/armies/move`,
                {
                    armySize: this.armySize,
                    fromCity: this.selectedTerritory.cityName,
                    toCity: this.data.cityName,
                },
                { headers }
            )
            .subscribe((response) => {
                if (response.status === true) this.dialogRef.close('');
            });
    }

    private _filterCities(value: string): any[] {
        const filterValue = value.toLowerCase();

        return this.ownedTerritory.filter((territory) =>
            territory.cityName.toLowerCase().includes(filterValue)
        );
    }
}

@Component({
    selector: 'city-name-change-dialog',
    templateUrl: './city-name-change-dialog.html',
})
export class CityNameChangeDialogComponent {
    constructor() {}
}

class March {
    fromCityName!: string;
    fromCityOwner!: string; // check if fromCityOwner or toCityOwner is you, then determine color
    toCityName!: string;
    toCityOwner!: string;
    returning!: boolean; // true if attacking another city and army is returning
    armySize!: number;
    startTime!: string;
    endTime!: string;
    attack!: boolean; // if you are attacking someone, then attack is true
    incoming!: boolean;
    text: string = '';
    ss!: number;
    dd!: number;
    hh!: number;
    mm!: number;
    progress!: number;
    color!: string;
}

@Component({
    selector: 'marches-dialog',
    templateUrl: './marches-dialog.html',
    styleUrls: ['./sidebar.component.css'],
})
export class MarchesDialogComponent {
    marches: March[] = [];
    cityOwner: string = '';

    marchTimeInterval!: ReturnType<typeof setInterval>;
    getMarchInterval!: ReturnType<typeof setInterval>;

    constructor(public http: HttpClient, private cookieService: CookieService) {
        let headers = new HttpHeaders();
        headers = headers.append('Token', this.cookieService.get('jwtToken'));
        this.http
            .get<any>(`${environment.API_HOST}/armies/marches`, { headers })
            .subscribe((response) => {
                this.marches = response ? response : [];
                for (let i = 0; i < this.marches.length; i++) {
                    let march = this.marches[i];
                    march.color = 'green';
                    if (march.returning) {
                        march.text = 'Returning';
                    } else if (march.attack) {
                        march.text = 'Attacking';
                        march.color = march.incoming ? 'red' : 'green';
                    } else {
                        march.text = 'Reinforcing';
                    }
                    if (march.startTime != '' && march.endTime != '') {
                        let startUnix = Date.parse(march.startTime);
                        let endUnix = Date.parse(march.endTime);
                        let time: number = Date.now();
                        march.progress =
                            ((time - startUnix) / (endUnix - startUnix)) * 100;

                        let remaining = endUnix - time;
                        march.ss = remaining / 1000;
                        march.dd = Math.floor(march.ss / 86400);
                        march.ss %= 86400;
                        march.hh = Math.floor(march.ss / 3600);
                        march.ss %= 3600;
                        march.mm = Math.floor(march.ss / 60);
                        march.ss %= 60;
                        march.ss = Math.floor(march.ss);
                    }
                }
            });

        this.marchTimeInterval = setInterval(() => {
            for (let i = 0; i < this.marches.length; i++) {
                let march = this.marches[i];
                if (march.startTime != '' && march.endTime != '') {
                    let startUnix = Date.parse(march.startTime);
                    let endUnix = Date.parse(march.endTime);
                    let time: number = Date.now();
                    march.progress =
                        ((time - startUnix) / (endUnix - startUnix)) * 100;

                    let remaining = endUnix - time;
                    march.ss = remaining / 1000;
                    march.dd = Math.floor(march.ss / 86400);
                    march.ss %= 86400;
                    march.hh = Math.floor(march.ss / 3600);
                    march.ss %= 3600;
                    march.mm = Math.floor(march.ss / 60);
                    march.ss %= 60;
                    march.ss = Math.floor(march.ss);
                }
            }
        }, 500);

        this.getMarchInterval = setInterval(() => {
            let headers = new HttpHeaders();
            headers = headers.append(
                'Token',
                this.cookieService.get('jwtToken')
            );
            this.http
                .get<any>(`${environment.API_HOST}/armies/marches`, { headers })
                .subscribe((response) => {
                    this.marches = response ? response : [];
                    for (let i = 0; i < this.marches.length; i++) {
                        let march = this.marches[i];
                        march.color = 'green';
                        if (march.returning) {
                            march.text = 'Returning';
                        } else if (march.attack) {
                            march.text = 'Attacking';
                            march.color = march.incoming ? 'red' : 'green';
                        } else {
                            march.text = 'Reinforcing';
                        }
                        if (march.startTime != '' && march.endTime != '') {
                            let startUnix = Date.parse(march.startTime);
                            let endUnix = Date.parse(march.endTime);
                            let time: number = Date.now();
                            march.progress =
                                ((time - startUnix) / (endUnix - startUnix)) *
                                100;

                            let remaining = endUnix - time;
                            march.ss = remaining / 1000;
                            march.dd = Math.floor(march.ss / 86400);
                            march.ss %= 86400;
                            march.hh = Math.floor(march.ss / 3600);
                            march.ss %= 3600;
                            march.mm = Math.floor(march.ss / 60);
                            march.ss %= 60;
                            march.ss = Math.floor(march.ss);
                        }
                    }
                });
        }, 2500);
    }

    public ngOnDestroy(): void {
        clearInterval(this.marchTimeInterval);
    }
}

@Component({
    selector: 'train-dialog',
    templateUrl: './train-dialog.html',
    styleUrls: ['./sidebar.component.css'],
})
export class TrainDialogComponent {
    constructor(
        public dialogRef: MatDialogRef<TrainDialogComponent>,
        public http: HttpClient,
        private cookieService: CookieService,
        private _snackBar: MatSnackBar,
        @Inject(MAT_DIALOG_DATA)
        public data: { cityName: string; maxTrainCount: number }
    ) {}
    maxCap: number = 0;
    trainArmySize: number = 1;
    armySize: number = 0;
    progBar: boolean = false;

    startTime!: string;
    endTime!: string;
    startUnix!: number;
    endUnix!: number;
    progress!: number;
    remaining!: number;
    dd!: number;
    hh!: number;
    mm!: number;
    ss!: number;

    interval!: ReturnType<typeof setInterval>;

    public ngOnInit(): void {
        this.maxCap = this.data.maxTrainCount;
        this.interval = setInterval(() => {
            let headers = new HttpHeaders();
            headers = headers.append(
                'Token',
                this.cookieService.get('jwtToken')
            );

            let parameter = '';
            let cityName = this.cookieService.get('cityName');
            if (cityName != '') {
                parameter = `?cityName=${encodeURIComponent(cityName)}`;
            }

            this.http
                .get<any>(
                    `${environment.API_HOST}/armies/training${parameter}`,
                    { headers }
                )
                .subscribe((response) => {
                    this.armySize = response.armySize;
                    this.startTime = response.startTime;
                    this.endTime = response.endTime;
                });

            if (this.startTime != '' && this.endTime != '') {
                this.progBar = true;
                this.startUnix = Date.parse(this.startTime);
                this.endUnix = Date.parse(this.endTime);
                if (
                    Number.isNaN(this.startUnix) ||
                    Number.isNaN(this.endUnix)
                ) {
                    this.progBar = false;
                }
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
        }, 200);
    }

    public train() {
        let headers = new HttpHeaders();
        headers = headers.append('Token', this.cookieService.get('jwtToken'));
        this.http
            .post<any>(
                `${environment.API_HOST}/armies/train`,
                {
                    troopCount: this.trainArmySize,
                    cityName: this.data.cityName,
                },
                { headers }
            )
            .subscribe((response) => {
                if (response.status) {
                    this.dialogRef.close('');
                } else {
                    this._snackBar.open('Error Training Troops!', 'Close', {
                        duration: 2000,
                    });
                }
            });
    }

    public ngOnDestroy(): void {
        clearInterval(this.interval);
    }
}

@Component({
    selector: 'delete-dialog',
    templateUrl: './delete-dialog.html',
    styleUrls: ['./sidebar.component.css'],
})
export class DeleteDialogComponent {
    constructor(
        public dialogRef: MatDialogRef<DeleteDialogComponent>,
        public http: HttpClient,
        private cookieService: CookieService,
        private _snackBar: MatSnackBar,
        @Inject(MAT_DIALOG_DATA)
        public data: { cityRow: number; cityColumn: number }
    ) {}

    public deleteBuilding() {
        let headers = new HttpHeaders();
        headers = headers.append('Token', this.cookieService.get('jwtToken'));
        let parameter = '';
        let cityName = this.cookieService.get('cityName');
        if (cityName != '') {
            parameter = `?cityName=${encodeURIComponent(cityName)}`;
        }

        this.http
            .post<any>(
                `${environment.API_HOST}/cities/destroyBuilding${parameter}`,
                {
                    cityRow: this.data.cityRow,
                    cityColumn: this.data.cityColumn,
                },
                { headers }
            )
            .subscribe((response) => {
                if (response.status) {
                    this.dialogRef.close('');
                } else {
                    this._snackBar.open(
                        'Error destroying the building!',
                        'Close',
                        {
                            duration: 2000,
                        }
                    );
                }
            });
    }
}

class Battle {
    fromCityName!: string;
    toCityName!: string;
    fromCityOwner!: string;
    toCityOwner!: string;
    attackerArmySize!: number;
    defenderArmySize!: number;
    amountLooted!: number;
    battleTime!: string;
    attackVictory!: boolean;
    incoming!: boolean;
    color!: string;
    hours!: number;
    days!: number;
    text!: string;
}

@Component({
    selector: 'battle-logs-dialog',
    templateUrl: './battle-logs-dialog.html',
    styleUrls: ['./sidebar.component.css'],
})
export class BattleLogsDialogComponent {
    battles: Battle[] = [];
    getBattlesInterval!: ReturnType<typeof setInterval>;

    constructor(public http: HttpClient, private cookieService: CookieService) {
        let headers = new HttpHeaders();
        headers = headers.append('Token', this.cookieService.get('jwtToken'));

        this.http
            .get<any>(`${environment.API_HOST}/armies/battles`, { headers })
            .subscribe((response) => {
                this.battles = response ? response : [];
                for (let i = 0; i < this.battles.length; i++) {
                    let battle = this.battles[i];

                    if (
                        (battle.incoming && !battle.attackVictory) ||
                        (!battle.incoming && battle.attackVictory)
                    ) {
                        battle.color = 'green';
                    } else {
                        battle.color = 'red';
                    }

                    battle.text =
                        battle.fromCityOwner + ' attacks ' + battle.toCityOwner;

                    battle.amountLooted = Math.round(battle.amountLooted);

                    if (battle.battleTime != '') {
                        let unix = Date.parse(battle.battleTime);
                        let remaining = Date.now() - unix;
                        let seconds = remaining / 1000;
                        battle.days = Math.floor(seconds / 86400);
                        seconds %= 86400;
                        battle.hours = Math.floor(seconds / 3600);
                    }
                }
            });

        this.getBattlesInterval = setInterval(() => {
            this.http
                .get<any>(`${environment.API_HOST}/armies/battles`, { headers })
                .subscribe((response) => {
                    this.battles = response ? response : [];
                    for (let i = 0; i < this.battles.length; i++) {
                        let battle = this.battles[i];
                        if (
                            (battle.incoming && !battle.attackVictory) ||
                            (!battle.incoming && battle.attackVictory)
                        ) {
                            battle.color = 'green';
                        } else {
                            battle.color = 'red';
                        }

                        battle.text =
                            battle.fromCityOwner +
                            ' attacks ' +
                            battle.toCityOwner;

                        battle.amountLooted = Math.round(battle.amountLooted);

                        if (battle.battleTime != '') {
                            let unix = Date.parse(battle.battleTime);
                            let remaining = Date.now() - unix;
                            let seconds = remaining / 1000;
                            battle.days = Math.floor(seconds / 86400);
                            seconds %= 86400;
                            battle.hours = Math.floor(seconds / 3600);
                        }
                    }
                });
        }, 1000);
    }

    public ngOnDestroy(): void {
        clearInterval(this.getBattlesInterval);
    }
}

class Territory {
    cityName!: string;
    cityProduction!: number;
    armySize!: number;
    cityPopulation!: number;
}

@Component({
    selector: 'territory-dialog',
    templateUrl: './territory-dialog.html',
    styleUrls: ['./sidebar.component.css'],
})
export class TerritoryDialogComponent {
    territory: Territory[] = [];

    constructor(
        public http: HttpClient,
        private cookieService: CookieService,
        public dialogRef: MatDialogRef<AttackDialogComponent>,
        private dialog: MatDialog
    ) {
        let headers = new HttpHeaders();
        headers = headers.append('Token', this.cookieService.get('jwtToken'));

        this.http
            .get<any>(`${environment.API_HOST}/cities/territory`, { headers })
            .subscribe((response) => {
                this.territory = response ? response : [];
            });
    }

    public visitCity(cityName: string): void {
        this.cookieService.set('cityName', cityName);
        this.dialogRef.close('');
    }

    public reinforce(cityName: string): void {
        this.dialog.open(MoveDialogComponent, {
            width: '1000px',
            height: '600px',
            data: { cityName: cityName },
        });
    }
}

@Component({
    selector: 'move-dialog',
    templateUrl: './move-dialog.html',
    styleUrls: ['./sidebar.component.css'],
})
export class MoveDialogComponent {
    ownedTerritory: CityArmy[] = [];
    territoryCtrl = new FormControl('');
    filteredTerritory!: Observable<CityArmy[]>;

    selectedTerritory!: CityArmy;
    selected: boolean = false;

    armySize: number = 1;

    constructor(
        public dialogRef: MatDialogRef<MoveDialogComponent>,
        public http: HttpClient,
        private cookieService: CookieService,
        @Inject(MAT_DIALOG_DATA) public data: { cityName: string }
    ) {
        let headers = new HttpHeaders();
        headers = headers.append('Token', this.cookieService.get('jwtToken'));
        this.http
            .get<any>(`${environment.API_HOST}/cities/territory`, { headers })
            .subscribe((response) => {
                this.ownedTerritory = response;
                for (let i = 0; i < this.ownedTerritory.length; i++) {
                    if (this.ownedTerritory[i].cityName == this.data.cityName) {
                        this.ownedTerritory.splice(i, 1);
                        break;
                    }
                }

                this.filteredTerritory = of(this.ownedTerritory);
                this.filteredTerritory = this.territoryCtrl.valueChanges.pipe(
                    startWith(''),
                    map((territory) =>
                        territory
                            ? this._filterCities(territory)
                            : this.ownedTerritory.slice()
                    )
                );
            });
    }

    public saveCityInfo(selectedTerritory: CityArmy) {
        this.selected = true;
        this.selectedTerritory = selectedTerritory;
    }

    public march() {
        let headers = new HttpHeaders();
        headers = headers.append('Token', this.cookieService.get('jwtToken'));
        this.http
            .post<any>(
                `${environment.API_HOST}/armies/move`,
                {
                    armySize: this.armySize,
                    fromCity: this.selectedTerritory.cityName,
                    toCity: this.data.cityName,
                },
                { headers }
            )
            .subscribe((response) => {
                if (response.status === true) this.dialogRef.close('');
            });
    }

    private _filterCities(value: string): any[] {
        const filterValue = value.toLowerCase();

        return this.ownedTerritory.filter((territory) =>
            territory.cityName.toLowerCase().includes(filterValue)
        );
    }
}
