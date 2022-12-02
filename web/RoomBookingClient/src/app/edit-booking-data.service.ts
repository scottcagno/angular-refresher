import { Injectable } from '@angular/core';
import {Room} from "./model/Room";
import {User} from "./model/User";
import {DataService} from "./data.service";

@Injectable({
  providedIn: 'root'
})
export class EditBookingDataService {

  rooms !: Array<Room>;
  users !: Array<User>;

  dataLoaded = 0;

  constructor(private dataService : DataService) {
    this.dataService.getRooms().subscribe(
      next => {
        this.rooms = next
        this.dataLoaded++
      }
    );
    this.dataService.getUsers().subscribe(
      next => {
        this.users = next
        this.dataLoaded++
      }
    );

  }
}
