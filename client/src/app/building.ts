export class Building {
    production_rate: number = 0.0;
    upgrade: number = 1;
    type: string = "";
    happiness_factor: number = 0.0;
    build_time: number = 0;
    level: number = 1;
    upgrade_cost: number = 0;
    building_icon: string = "";
    // production rate
    // upgrade
    // type
    // happiness factor
    // build time
    // building level
    // upgrade cost
    // building icon

    // building types
    
    constructor(type: string) {
        if (type == "city_hall") {
            this.production_rate = 0;
            this.type = type;
            this.happiness_factor = 0.0;
            this.build_time = 0;
            this.upgrade_cost = 1000;
            this.building_icon = "🕋";
        }
    }
}
