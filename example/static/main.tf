provider "gobetween" {
    host = "35.189.10.204"
    port = 8888
}

resource "gobetween_server" "example" {
    name = "example"
    balance = "weight"
    bind = "0.0.0.0:5858"
    discovery {
        static_list = ["1.2.3.4:80", "2.3.4.5:80"]
    }
    # static_backends = ["1.2.3.4:80", "2.3.4.5:80"]
}