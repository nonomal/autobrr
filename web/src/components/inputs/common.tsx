import { Field } from "formik";

interface ErrorFieldProps {
    name: string;
    classNames?: string;
    subscribe?: any;
}

const ErrorField = ({ name, classNames }: ErrorFieldProps) => (
    <Field name={name} subscribe={{ touched: true, error: true }}>
        {({ meta: { touched, error } }: any) =>
            touched && error ? <span className={classNames}>{error}</span> : null
        }
    </Field>
);

interface CheckboxFieldProps {
    name: string;
    label: string;
    sublabel?: string;
}

const CheckboxField = ({
    name,
    label,
    sublabel
}: CheckboxFieldProps) => (
    <div className="relative flex items-start">
        <div className="flex items-center h-5">
            <Field  
                id={name}
                name={name}
                type="checkbox" 
                className="focus:ring-bkue-500 h-4 w-4 text-blue-600 border-gray-300 rounded"
            />
        </div>
        <div className="ml-3 text-sm">
        <label htmlFor={name} className="font-medium text-gray-900 dark:text-gray-100">
            {label}
        </label>
        <p className="text-gray-500">{sublabel}</p>
        </div>
    </div>
)

export { ErrorField, CheckboxField }