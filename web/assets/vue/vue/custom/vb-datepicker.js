Vue.component('date-picker', {
    template: '<input/>',
    props: [
        'config', 'value',
    ],
    data: function() {
        return {
            configDefault: {
                autoclose: true,
                format: 'dd/mm/yyyy',
                todayBtn: 'linked',
                clearBtn: true,
            }
        };
    },
    mounted: function () {
        var self = this;

        $(this.$el).datepicker(this.config != undefined ? this.config : this.configDefault)
            .on('changeDate', function (e) {
                self.$emit('update-date', e);
                self.$emit('input', e.format('yyyy-mm-dd'));
            });

        if (this.value != '' && this.value != null) {
            date = new Date(this.value);
            $(this.$el).datepicker("setDate", date);
            $(this.$el).datepicker('update');
        }
    },
    beforeDestroy: function () {
        $(this.$el).datepicker('hide').datepicker('destroy');
    }
});