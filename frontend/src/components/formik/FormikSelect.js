import React from "react";
import { FormControl, InputLabel, Select } from "@mui/material";
import { useField } from "formik";
import PropTypes from "prop-types";

const FormikSelect = (props) => {
    const [field] = useField(props.name);

    const {onChange} = field;

    const {className, options, id, ...otherProps} = props;

    return (
        <FormControl className={className}>
            <InputLabel id={id}>Status</InputLabel>
            <Select
                native
                labelId={id}
                onChange={onChange}
                {...otherProps}
            >
                {options.map((option) => <option key={option.value} value={option.value}>{option.display}</option>)}
            </Select>
        </FormControl>
    );
};

FormikSelect.propTypes = {
    name: PropTypes.string.isRequired,
    id: PropTypes.string.isRequired,
    options: PropTypes.array.isRequired
};

export default FormikSelect;
