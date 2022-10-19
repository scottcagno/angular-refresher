import {Component, OnDestroy, OnInit} from '@angular/core';
import {DataService} from "../../data.service";
import {Room} from "../../model/Room";
import {ActivatedRoute, Router} from "@angular/router";

@Component({
  selector: 'app-rooms',
  templateUrl: './rooms.component.html',
  styleUrls: ['./rooms.component.css']
})
export class RoomsComponent implements OnInit {

  rooms !:Array<Room>;
  selectedRoom !:Room;
  action !:string;

  constructor(private dataService :DataService, private route :ActivatedRoute, private router :Router) {}

  ngOnInit(): void {
    this.dataService.getRooms().subscribe(
      next => {
        this.rooms = next
      }
    );
    this.route.queryParams.subscribe(
      (params) => {
        const id = params['id'];
        if (id) {
          // @ts-ignore
          this.selectedRoom = this.rooms.find((room)=> {
            // + symbol before a string, converts a string into a number
            return room.id === +id
          });
          this.action = params['action'];
        }
        if (params['action'] == 'add') {
          this.selectedRoom = new Room();
          this.action = 'edit';
        }
      }
    );
  }

  selectRoom(id :number) {
    this.router.navigate(['admin','rooms'], {
      queryParams:{id:id, action: 'view'}
    })
  }

  addRoom() {
    this.router.navigate(['admin','rooms'], {
      queryParams:{action: 'add'}
    })
  }

}
