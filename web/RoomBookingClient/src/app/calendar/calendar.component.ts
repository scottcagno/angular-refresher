import { Component, OnInit } from '@angular/core';
import {formatDate} from "@angular/common";
import {Booking} from "../model/Booking";
import {DataService} from "../data.service";
import {ActivatedRoute, Router} from "@angular/router";

@Component({
  selector: 'app-calendar',
  templateUrl: './calendar.component.html',
  styleUrls: ['./calendar.component.css']
})
export class CalendarComponent implements OnInit {

  bookings !: Array<Booking>;
  selectedDate !: string;

  constructor(private dataService: DataService,
              private router :Router,
              private route :ActivatedRoute) {
  }

  ngOnInit(): void {
    this.route.queryParams.subscribe(params =>{
      this.selectedDate = params['date'];
      if (!this.selectedDate) {
        this.selectedDate = formatDate(new Date(), 'yyyy-MM-dd', 'en_us');
      }
      this.dataService.getBookings(this.selectedDate).subscribe(
        next => {
          this.bookings = next;
        }
      );
    });
  }

  editBooking(id: number) {
    this.router.navigate(['edit-booking'],{queryParams:{id:id}});
  }

  addBooking() {
    this.router.navigate(['add-booking']);
  }

  deleteBooking(id :number) {
    this.dataService.deleteBooking(id).subscribe();
  }

  dateChanged() {
    this.router.navigate([''], {queryParams:{date:this.selectedDate}})
  }

}

const dateFormats = [
"<hr>",
"<p>{{ selectedDate | date }}</p>",
"<p>{{ selectedDate | date:'yyyy-MM-dd' }}</p>",
  "<p>{{ selectedDate | date:'MMM dd yy' }}</p>",
  "<p>{{ selectedDate | date:'MMMM YY' }}</p>",
  "<hr>",
  "<p>{{ selectedDate | date: 'short' }}</p>",
  "<p>{{ selectedDate | date: 'medium' }}</p>",
  "<p>{{ selectedDate | date: 'long' }}</p>",
  "<p>{{ selectedDate | date: 'full' }} </p>",
  "<p>{{ selectedDate | date: 'shortDate' }}</p>",
  "<p>{{ selectedDate | date: 'mediumDate' }}</p>",
  "<p>{{ selectedDate | date: 'longDate' }}</p>",
  "<p>{{ selectedDate | date: 'fullDate' }}</p>",
  "<p>{{ selectedDate | date: 'shortTime' }}</p>",
  "<p>{{ selectedDate | date: 'mediumTime' }}</p>",
  "<p>{{ selectedDate | date: 'longTime' }}</p>",
  "<p>{{ selectedDate | date: 'fullTime' }}</p>",
  "<hr>",
  "<p>{{ selectedDate | date: 'M/d/yy, h:mm a' }}</p>",
  "<p>{{ selectedDate | date: 'MMM d, y, h:mm:ss a' }}</p>",
  "<p>{{ selectedDate | date: 'MMMM d, y, h:mm:ss a z' }}</p>",
  "<p>{{ selectedDate | date: 'EEEE, MMMM d, y, h:mm:ss a zzzz' }} </p>",
  "<p>{{ selectedDate | date: 'M/d/yy' }}</p>",
  "<p>{{ selectedDate | date: 'MMM d, y' }}</p>",
  "<p>{{ selectedDate | date: 'MMMM d, y' }}</p>",
  "<p>{{ selectedDate | date: 'EEEE, MMMM d, y' }}</p>",
  "<p>{{ selectedDate | date: 'h:mm a' }}</p>",
  "<p>{{ selectedDate | date: 'h:mm:ss a' }}</p>",
  "<p>{{ selectedDate | date: 'h:mm:ss a z' }}</p>",
  "<p>{{ selectedDate | date: 'h:mm:ss a zzzz' }}</p>",
];
