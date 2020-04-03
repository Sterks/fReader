<template>
    <form action="/tasks" method="POST">
        <p>Check some colors</p>
        <div class="row">
            <div class="input-field col s12 m3">
                <input type="text" id="ts_name" v-model="ts_name">
                <label for="ts_name">Название задачи</label>
            </div>
        </div>
        <div class="row">
            <div class="input-field col s12 m3">
                <input type="text" id="ts_data_start" v-model="ts_data_start">
                <label for="ts_data_start">Дата создания</label>
            </div>
        </div>
        <div class="row">
            <div class="input-field col s12 m3">
                <input type="text" class="datepicker" id="ts_run_times" v-model="ts_run_times">
                <label for="ts_run_times">Кол-во раз</label>
            </div>
        </div>
        <div class="row">
            <div class="input-field col s12 m3">
                <input type="text" id="ts_comment" v-model="ts_comment">
                <label for="ts_comment">Комментарий</label>
            </div>
        </div>

        <button class="btn btn-success">Отправить</button>

    </form>
</template>

<script>
    import M from 'materialize-css/dist/js/materialize.js'
    document.addEventListener('DOMContentLoaded', function() {
        document.querySelectorAll('.datepicker');
        M.Datepicker.init();
    });

    export default {
        data() {
            return {

            }
        }
    }
</script>