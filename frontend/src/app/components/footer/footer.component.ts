import { Component, OnInit, ElementRef, Renderer2, HostListener } from '@angular/core';
import { TxnTokenService } from '../../services/txn-token.service';

@Component({
  selector: 'app-footer',
  templateUrl: './footer.component.html',
  styleUrls: ['./footer.component.css'],
  standalone: false
})
export class FooterComponent implements OnInit {
  tratJwt: string | null = null;
  trat: string | null = null;
  isDragging: boolean = false;
  startY: number = 0;
  startHeight: number = 0;
  minHeight: number = 10;
  maxHeight: number = 50;

  constructor(
    private txnTokenService: TxnTokenService,
    private el: ElementRef,
    private renderer: Renderer2
  ) { }

  ngOnInit(): void {
    this.txnTokenService.txnToken$.subscribe(token => {
      this.tratJwt = token;
      this.trat = this.decodeJwtPayload(token);
      this.updateFooterVisibility();
    });
  }

  private decodeJwtPayload(token: string | null): string | null {
    if (!token) return null;
    try {
      const base64Payload = token.split('.')[1];
      const payload = JSON.parse(atob(base64Payload));
      return JSON.stringify(payload, null, 2);
    } catch (e) {
      console.error('Error decoding JWT payload:', e);
      return null;
    }
  }

  private updateFooterVisibility(): void {
    if (this.tratJwt) {
      this.renderer.setStyle(this.el.nativeElement, 'display', 'flex');
    } else {
      this.renderer.setStyle(this.el.nativeElement, 'display', 'none');
    }
  }

  @HostListener('mousedown', ['$event'])
  onMouseDown(event: MouseEvent) {
    if (event.offsetY < 10) {
      this.isDragging = true;
      this.startY = event.clientY;
      this.startHeight = this.el.nativeElement.offsetHeight;
      this.addGlobalListeners();
    }
  }

  @HostListener('document:mousemove', ['$event'])
  onMouseMove(event: MouseEvent) {
    if (this.isDragging) {
      const diffY = this.startY - event.clientY;
      let newHeight = (this.startHeight + diffY) / window.innerHeight * 100;
      newHeight = Math.max(this.minHeight, Math.min(this.maxHeight, newHeight));
      this.renderer.setStyle(this.el.nativeElement, 'height', `${newHeight}vh`);
      event.preventDefault();
    }
  }

  @HostListener('document:mouseup')
  onMouseUp() {
    if (this.isDragging) {
      this.isDragging = false;
      this.removeGlobalListeners();
    }
  }

  private addGlobalListeners() {
    this.renderer.addClass(document.body, 'resize-active');
  }

  private removeGlobalListeners() {
    this.renderer.removeClass(document.body, 'resize-active');
  }
}