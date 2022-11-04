import {Component, OnDestroy, OnInit} from '@angular/core';
import {DataService} from "../../data.service";
import {Room} from "../../model/Room";
import {ActivatedRoute, Router} from "@angular/router";
import {FormResetService} from "../../form-reset.service";

@Component({
  selector: 'app-rooms',
  templateUrl: './rooms.component.html',
  styleUrls: ['./rooms.component.css']
})
export class RoomsComponent implements OnInit {

  rooms !:Array<Room>;
  selectedRoom !:Room;
  action !:string;
  loadingData = true;
  message = "One moment please... getting the list of rooms."
  maxRetry = 10;

  constructor(private dataService :DataService,
              private route :ActivatedRoute,
              private router :Router,
              private formResetService: FormResetService) {}


  loadData() {
    this.dataService.getRooms().subscribe(
      next => {
        this.rooms = next;
        this.loadingData = false;
        this.processURLParams();
      },
      (error) => {
        if (error.status === 0) {
          this.maxRetry--;
          if (this.maxRetry > 0) {
            this.message = "Sorry, something went wrong, trying again, please wait."
            this.loadData();
          } else {
            this.message = `Whoops, we seem to have encountered a ${error.status} error.`;
          }
        }
      }
    );
  }

  processURLParams() {
    this.route.queryParams.subscribe(
      (params) => {
        this.action = '';
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
          this.formResetService.resetRoomFormEvent.emit(this.selectedRoom);
        }
      }
    );
  }

  ngOnInit(): void {

    this.loadData();


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
