import { Injectable } from '@angular/core';
import {Layout, LayoutCapacity, Room} from "./model/Room";

@Injectable({
  providedIn: 'root'
})
export class DataService {

  rooms !: Array<Room>;

  constructor() {
    this.rooms = new Array<Room>();

    const room1 = new Room(1, 'first room', 'first floor');
    const cap1 = new LayoutCapacity(Layout.THEATER, 50);
    const cap2 = new LayoutCapacity(Layout.USHAPE, 20);
    room1.capacities.push(cap1, cap2);

    const room2 = new Room(2, 'second room', 'third floor');
    const cap3 = new LayoutCapacity(Layout.THEATER, 60);
    room2.capacities.push(cap3);

    this.rooms.push(room1, room2);
  }
}
