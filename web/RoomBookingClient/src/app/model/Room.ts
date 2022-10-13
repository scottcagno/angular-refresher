export class Room {

  id !:number;
  name !:string;
  location !:string;
  capacities !:Array<LayoutCapacity>;

  constructor(name ?:string, location ?:string) {
    if (name) { this.name = name }
    if (location) { this.location = location }
  }

}

export class LayoutCapacity {
  layout !:Layout;
  capacity !:number;
}

export enum Layout {
  THEATER = 'Theater',
  USHAPE = 'U-Shape',
  BOARD = 'Board Meeting'
}
