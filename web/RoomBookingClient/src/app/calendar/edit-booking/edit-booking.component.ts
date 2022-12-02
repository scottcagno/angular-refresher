import { Component, OnInit } from '@angular/core';
import {Booking} from "../../model/Booking";
import {Layout, Room} from "../../model/Room";
import {DataService} from "../../data.service";
import {User} from "../../model/User";
import {ActivatedRoute, Router} from "@angular/router";
import {map} from "rxjs";

@Component({
  selector: 'app-edit-booking',
  templateUrl: './edit-booking.component.html',
  styleUrls: ['./edit-booking.component.css']
})
export class EditBookingComponent implements OnInit {


  booking !: Booking;
  rooms !:Array<Room>;
  layouts = Object.keys(Layout);
  layoutEnum = Object(Layout);
  users !:Array<User>;
  dataLoaded = false;
  message = 'Please waitzzzzzzzzzzzzzz...';


  constructor(private dataService :DataService,
              private route :ActivatedRoute,
              private router :Router) { }

  ngOnInit(): void {

    this.rooms = this.route.snapshot.data['rooms']
    this.users = this.route.snapshot.data['users']

    const id = this.route.snapshot.queryParams['id'] as number;
    if (id) {
      this.dataService.getBooking(+id)
        .pipe(
          map(booking => {
            booking.room = this.rooms.find(room => room.id === booking.room.id) as Room;
            booking.user = this.users.find(user => user.id === booking.user.id) as User;
           this.dataLoaded = true;
           booking.id = id;
            return booking;
          })
        )
        .subscribe(
        next => {
          this.booking = next;
          this.dataLoaded = true;
          this.message = '';
        },
      );
    } else {
      this.booking = new Booking({});
      this.dataLoaded = true;
      this.message = '';
    }
  }

  onCancel() {
    this.router.navigate(['']);
  }

  onSubmit() {
    if (this.booking.id != null) {
      this.dataService.saveBooking(this.booking).subscribe(
        next => {
          this.router.navigate(['']);
        }
      );
    } else {
      this.dataService.addBooking(this.booking).subscribe(
        next => {
          this.router.navigate(['']);
        }
      );
    }
  }

}
