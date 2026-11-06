'use client';
import { FormDescription, FormLabel } from '@/components/ui/form';

import FormErrorMessage from '@/components/FormErrorMessage';
import { Switch } from '@/components/ui/switch';
import { create } from '@bufbuild/protobuf';
import { GenerateCpf, GenerateCpfSchema } from '@neosync/sdk';
import { ReactElement } from 'react';
import { TransformerConfigProps } from './util';

interface Props extends TransformerConfigProps<GenerateCpf> {}

export default function GenerateCpfForm(props: Props): ReactElement {
  const { value, setValue, isDisabled, errors } = props;

  return (
    <div className="flex flex-row items-center justify-between rounded-lg border dark:border-gray-700 p-3 shadow-xs">
      <div className="space-y-0.5 w-[80%]">
        <FormLabel>Formatted</FormLabel>
        <FormDescription>
          Format the CPF as xxx.xxx.xxx-xx. If disabled, only the 11 digits are
          returned.
        </FormDescription>
      </div>
      <div className="flex flex-col">
        <div className="justify-end flex">
          <Switch
            checked={value.formatted ?? true}
            onCheckedChange={(checked) => {
              setValue(
                create(GenerateCpfSchema, {
                  ...value,
                  formatted: checked,
                })
              );
            }}
            disabled={isDisabled}
          />
          <FormErrorMessage message={errors?.formatted?.message} />
        </div>
      </div>
    </div>
  );
}
