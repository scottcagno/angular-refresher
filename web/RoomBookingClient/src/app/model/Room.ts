export class Room {

  id !:number;
  name !:string;
  location !:string;
  capacities  !:Array<LayoutCapacity>;

  constructor(id ?:number, name ?:string, location ?:string) {
    if (id) { this.id = id }
    if (name) { this.name = name }
    if (location) { this.location = location }
    this.capacities = new Array<LayoutCapacity>();
  }

}

export class LayoutCapacity {
  layout !:Layout;
  capacity !:number;

  constructor(layout ?: Layout, capacity ?: number) {
    if (layout) { this.layout = layout }
    if (capacity) { this.capacity = capacity }
  }
}

export enum Layout {
  THEATER = 'Theater',
  USHAPE = 'U-Shape',
  BOARD = 'Board Meeting'
}
