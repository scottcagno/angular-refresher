import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {Room} from "../../../model/Room";
import {Router} from "@angular/router";
import {DataService} from "../../../data.service";
import {ConfirmDialogService} from "../../../confirm-dialog/confirm-dialog-service";
import { NgbModal } from '@ng-bootstrap/ng-bootstrap';
import {AuthService} from "../../../auth.service";

@Component({
  selector: 'app-room-detail',
  templateUrl: './room-detail.component.html',
  styleUrls: ['./room-detail.component.css']
})
export class RoomDetailComponent implements OnInit {

  @Input()
  room !:Room;

  @Output()
  dataChanged = new EventEmitter();

  message = '';
  isAdminUser = false;

  constructor(
    private router :Router,
    private dataService : DataService,
    private authService: AuthService) { }

  ngOnInit(): void {
    if (this.authService.getRole() === 'ROLE_ADMIN') {
      this.isAdminUser = true;
    }
  }

  editRoom() {
    this.router.navigate(['admin','rooms'], {queryParams:{id:this.room.id, action: 'edit'}});
  }

  deleteRoom() {
    const result = confirm('Are you sure you wish to delete this room?');
    if (result) {
      this.message = 'Removing room';
      this.dataService.deleteRoom(this.room.id).subscribe(
        next => {
          this.dataChanged.emit();
          this.router.navigate(['admin','rooms']);
        },
        error => {
          this.message = 'Sorry -- this room cannot be deleted at this time.';
        }
      );
    }
  }

}
