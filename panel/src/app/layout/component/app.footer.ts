import { Component } from '@angular/core';

@Component({
    standalone: true,
    selector: 'app-footer',
    template: `<div class="layout-footer">
        GENEZ by
        <a href="https://mstgnz.com" target="_blank" rel="noopener noreferrer" class="text-primary font-bold hover:underline">GENEZ</a>
    </div>`
})
export class AppFooter {}
