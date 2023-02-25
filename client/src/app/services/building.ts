export class Building {
    type: string = "";
    level: number = 1;
    name: string = "";
    production: number = 0.0;
    happiness_change: number = 0.0;
    start_time: string = "";
    end_time: string = "";
    building_icon: string = "";
    build_cost: number = 0.0;
    build_time: number = 0;
    building_icons: Map<string, string> = new Map([
        ["City Hall", "🏛"],
        ["Hospital", "🏥"]
    ]);
    row: number = 0;
    column: number = 0;
    // production rate
    // upgrade
    // type
    // happiness factor
    // build time
    // building level
    // upgrade cost
    // building icon

    // building types
    
    constructor(level: number = 0, name: string = "", type: string = "", row: number = 0, column: number = 0) {
        this.level = level;
        this.name = name;
        this.type = type;
        this.row = row;
        this.column = column;
        this.building_icon = this.building_icons.get(type) as string;
        // initialize stats
    }

    setIcon() {
        this.building_icon = this.building_icons.get(this.type) as string;
    }
}