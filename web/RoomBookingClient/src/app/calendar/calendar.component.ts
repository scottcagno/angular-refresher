import { Component, OnInit } from '@angular/core';
import {formatDate} from "@angular/common";

@Component({
  selector: 'app-calendar',
  templateUrl: './calendar.component.html',
  styleUrls: ['./calendar.component.css']
})
export class CalendarComponent implements OnInit {

  constructor() { }

  selectedDate = new Date();

  ngOnInit(): void {
    const date = formatDate(this.selectedDate, 'medium', 'en_us');
    console.log(date);
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
