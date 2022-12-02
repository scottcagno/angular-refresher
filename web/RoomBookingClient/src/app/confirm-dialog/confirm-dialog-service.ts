import { Injectable } from '@angular/core';
import { NgbModal } from '@ng-bootstrap/ng-bootstrap';
import {ConfirmDialogComponent} from "./confirm-dialog.component";

@Injectable({
  providedIn: 'root'
})
export class ConfirmDialogService {

  constructor(private modalService: NgbModal) {}

  public confirm(
    title: string,
    message: string,
    btnAcceptText: string = 'OK',
    btnCancelText: string = 'Cancel',
    dialogSize: 'sm'|'lg' = 'sm'): Promise<boolean> {


    const modalRef = this.modalService.open(ConfirmDialogComponent, {size: dialogSize});
    modalRef.componentInstance.title = title;
    modalRef.componentInstance.message = message;
    modalRef.componentInstance.btnAcceptText = btnAcceptText;
    modalRef.componentInstance.btnCancelText = btnCancelText;

    return modalRef.result;
  }
}

