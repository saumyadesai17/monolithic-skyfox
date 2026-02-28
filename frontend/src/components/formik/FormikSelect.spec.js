import React from "react";
import {FormikSelect} from "./index";
import * as formik from "formik";
import { render, screen } from "@testing-library/react";

// Mock the MUI components
vi.mock("@mui/material", () => ({
    FormControl: ({ children, className }) => <div data-testid="form-control" className={className}>{children}</div>,
    InputLabel: ({ children, id }) => <label data-testid="input-label" id={id}>{children}</label>,
    Select: ({ children, native, labelId, onChange, name, ...props }) => (
        <select 
            data-testid="select" 
            data-native={native ? "true" : "false"}
            data-labelid={labelId}
            name={name}
            onChange={onChange}
            {...props}
        >
            {children}
        </select>
    )
}));

vi.mock("formik");

describe("Basic Rendering", () => {
    let field;

    beforeEach(() => {
        field = {
            value: "test value",
            onChange: vi.fn(),
            onBlur: vi.fn()
        };
    });

    it("Should render with correct options", () => {
        // Mock useField to return our test data
        formik.useField = vi.fn().mockReturnValue([field, {}]);

        render(
            <FormikSelect 
                testProp='test prop' 
                name="test select" 
                id="test id"
                options={[
                    {value: "valueOne", display: "Value One"},
                    {value: "valueTwo", display: "Value Two"}
                ]}
            />
        );

        // Check that the input label is rendered with the correct ID
        const inputLabel = screen.getByTestId("input-label");
        expect(inputLabel).toHaveAttribute("id", "test id");
        expect(inputLabel).toHaveTextContent("Status");

        // Check that the select is rendered with the correct attributes
        const select = screen.getByTestId("select");
        expect(select).toHaveAttribute("data-native", "true");
        expect(select).toHaveAttribute("data-labelid", "test id");
        expect(select).toHaveAttribute("name", "test select");
        expect(select).toHaveAttribute("testProp", "test prop");

        // Check that the options are rendered correctly
        const options = screen.getAllByRole("option");
        expect(options).toHaveLength(2);
        expect(options[0]).toHaveValue("valueOne");
        expect(options[0]).toHaveTextContent("Value One");
        expect(options[1]).toHaveValue("valueTwo");
        expect(options[1]).toHaveTextContent("Value Two");
    });
});
