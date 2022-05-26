Vue.component('range-date-picker', {
    template: '<input/>',
    props: [
        'config', 'value'
    ],
    mounted: function () {
        var self = this;
        $(this.$el).daterangepicker(this.config)
            .on('apply.daterangepicker', function (ev, picker) {
                $(this).val(picker.startDate.format('MM/DD/YYYY') + ' - ' + picker.endDate.format('MM/DD/YYYY'));
                self.$emit('apply-date', picker);
            }).on('cancel.daterangepicker', function (ev, picker) {
                $(this).val('');
                self.$emit('cancel-date', picker);
            });
    },
    beforeDestroy: function () {
        $(this.$el).daterangepicker('destroy')
    }
});