class InputField {
        constructor() {
                this.options = ["Art Deco", "Colonial/Neo-Colonial", "Neoclassical/Greek Revival"]
                this.value = "Art Deco"
                this.dropDown = false;
        }

        switchInput() {

        }

        toggleDropDown(display) {
                display === "none" ? this.openDropDown() : this.closeDropDown();
        }

        closeDropDown() {
                document.getElementById("style-select").style.display = "none"
        }

        openDropDown() {
                const element = document.getElementById("style-select")
                element.style.display = "block"
        }

        changeOption(value) {
                this.value = value;
                const select = document.getElementById("select");
                select.value = value;
                this.closeDropDown()
        }

        async sendInput() {
                if (!this.value || this.value === "") return;
                const response = await fetch("/api/hello");
                console.log(response)
        }
}

const input = new InputField();
const select = document.getElementById("select");
select.addEventListener("click", () => {
        const display = getComputedStyle(document.getElementById("style-select")).getPropertyValue('display')
        input.toggleDropDown(display);
})

const arr = [...document.getElementsByClassName("option")]
arr.forEach((node) => {
        node.addEventListener("click", () => {
                input.changeOption(node.getAttribute("data-value"))
        });
})

document.getElementById("submit").addEventListener("click", () => {
        input.sendInput();
})