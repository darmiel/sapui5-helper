package gen

import "fmt"

func GenerateXMLView(namespace, view string) string {
	return fmt.Sprintf(`<mvc:View xmlns:mvc="sap.ui.core.mvc" xmlns="sap.m" controllerName="%s.controller.%s" displayBlock="true"><sap.ui.layout:VerticalLayout xmlns:sap.ui.layout="sap.ui.layout" width="100%%" id="layout0">
    <sap.ui.layout:content>
			<!-- TODO: Add contents here -->
			<Label text="Hello World!" width="100%%" id="label0" />
    </sap.ui.layout:content>
    </sap.ui.layout:VerticalLayout>
</mvc:View>`, namespace, view)
}

func GenerateJSView(namespace, view string) string {
	return fmt.Sprintf(`sap.ui.jsview("%s.view.%s", {

    getControllerName: function () {
        return "%s.controller.%s";
    },

    createContent: function (ctl) {
		// TODO: Implement me

        const eText = new sap.m.Text({
            text: "Hello World!"
        });
        return [eText];
    }
});`, namespace, view, namespace, view)
}

func GenerateController(namespace, view string) string {
	return fmt.Sprintf(`sap.ui.define([
    "sap/ui/core/mvc/Controller"
],
	/**
	 * @param {typeof sap.ui.core.mvc.Controller} Controller
	 */
    function (Controller) {
        "use strict";

        return Controller.extend("%s.controller.%s", {
			// Add code here
		});
	});`, namespace, view)
}
