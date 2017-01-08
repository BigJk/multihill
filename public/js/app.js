main = new Vue({
    el: 'body',
    data: {
        Hills: [],
        Hill: null,
        Selected: "",
        Submit: false
    },
    ready: function() {
        this.$http.get('hills').success(function(hillNames) {
            this.$set('Hills', hillNames.sort());
        }).error(function(error) {
            console.log(error);
        });
    },
    computed: {
        hillCount: function() {
            return this.Hills.length;
        }
    },
    methods: {
        select: function(hill) {
            this.closeSubmit();
            addClass("hill-" + hill, "active");
            if (this.Selected != "") {
                removeClass("hill-" + this.Selected, "active")
            }
            this.$set('Selected', hill)
            this.refresh();
        }, refresh: function() {
            this.$http.get('hill/' + this.Selected).success(function(hill) {
                this.$set('Hill', hill);
            }).error(function(error) {
                console.log(error);
            });
        }, closeSubmit: function() {
            this.$set('Submit', false);
        }, openSubmit: function() {
            this.$set('Submit', true);
        }, submit: function() {
            this.$http.post('hill/' + this.Selected, { Code: document.getElementById("warrior-code").value }).success(function(msg) {
                if (msg != "Ok") {
                    document.getElementById("error-msg").innerHTML = "Error: " + msg;
                } else {
                    this.closeSubmit();
                    this.refresh();
                }
            }).error(function(error) {
                console.log(error);
            });
        }, upload: function() {
            if (!window.FileReader) {
                alert('Your browser is not supported');
                return false;
            }
            var fileInput = document.getElementById("file-upload");
            for (var i = 0; i < fileInput.files.length; i++) {
                var textFile = fileInput.files[i];
                var reader = new FileReader();
                reader.readAsText(textFile);
                reader.addEventListener("load", this.uploadFile);
            }
        }, uploadFile: function(e) {
            var c = e.target.result;
            this.$http.post('hill/' + this.Selected, { Code: c }).success(function(msg) {
                if (msg == "Ok") {
                    this.closeSubmit();
                    this.refresh();
                }
            }).error(function(error) {
                console.log(error);
            });
        }
    }
});

setInterval(main.refresh, 2000);