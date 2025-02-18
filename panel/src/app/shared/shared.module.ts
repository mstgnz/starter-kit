import { NgModule } from '@angular/core';
import { DetePipeComponent } from './pipes/date.pipe';
import { PhonePipeComponent } from './pipes/phone.pipe';
import { RangePipeComponent } from './pipes/range.pipe';
import { StringPipeComponent } from './pipes/string.pipe';
import { PercentPipeComponent } from './pipes/percent.pipe';
import { ExcelExportComponent } from './excel/excel-export.component';
import { ExcelImportComponent } from './excel/excel-import.component';

@NgModule({
    declarations: [

    ],
    providers: [],
    imports: [
        DetePipeComponent,
        RangePipeComponent,
        PhonePipeComponent,
        StringPipeComponent,
        PercentPipeComponent,
        ExcelExportComponent,
        ExcelImportComponent
    ],
    exports: [
        DetePipeComponent,
        RangePipeComponent,
        PhonePipeComponent,
        StringPipeComponent,
        PercentPipeComponent,
        ExcelExportComponent,
        ExcelImportComponent,
    ]
})
export class SharedModule { }