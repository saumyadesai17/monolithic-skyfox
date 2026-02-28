import React from "react";
import {FormikTextField} from ".";
import { render, screen } from "@testing-library/react";
import * as formik from "formik";

vi.mock("formik");

describe("Basic rendering", () => {
    let field;
    let meta;

    beforeEach(() => {
        field = {
            value: "test value field",
            onChange: vi.fn(),
            onBlur: vi.fn()
        };

        meta = {
            error: "test error",
            touched: true
        };
    });

    function basicAssertions() {
        // In React Testing Library, we test what the user sees rather than implementation details
        expect(screen.getByLabelText("test label")).toBeInTheDocument();
        expect(screen.getByLabelText("test label")).toHaveValue("test value field");
    }

    it("Should render a formik text field correctly with errors", () => {
        // Mock useField to return our test data
        formik.useField = vi.fn().mockReturnValue([field, meta]);

        render(<FormikTextField testProp="test prop value" name="test field" label="test label"/>);
        
        basicAssertions();
        expect(screen.getByText("test error")).toBeInTheDocument();
    });

    it("Should render a formik text field correctly without errors", () => {
        meta.touched = false;
        meta.error = 'error text';
        
        // Mock useField to return our test data
        formik.useField = vi.fn().mockReturnValue([field, meta]);

        render(<FormikTextField testProp="test prop value" name="test field" label="test label"/>);

        basicAssertions();
        // When not touched, no error should be displayed
        expect(screen.queryByText("error text")).not.toBeInTheDocument();
    });
});
