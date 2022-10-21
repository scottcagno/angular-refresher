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
