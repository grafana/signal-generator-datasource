import React, { PureComponent } from 'react';
import { Button, Icon, InlineField, Input, Select, UnitPicker } from '@grafana/ui';
import { SignalField } from '../types';
import { formatLabels, parseLabels, SelectableValue } from '@grafana/data';
import { standardWaves } from 'info';

interface Props {
  signal: SignalField;
  index: number;
  onChange: (value: SignalField | undefined, index: number) => void;
  onAddExpr?: () => void;
}

export class SignalFieldEditor extends PureComponent<Props> {
  onExprInputChanged = (v: React.SyntheticEvent<HTMLInputElement>) => {
    this.onNewExpr(v.currentTarget.value);
  };

  onExprChange = (sel: SelectableValue<string>) => {
    this.onNewExpr(sel.value!);
  };

  onNewExpr = (expr: string) => {
    const { onChange, signal, index } = this.props;
    onChange({ ...signal, expr }, index);
  };

  onNameChange = (v: React.SyntheticEvent<HTMLInputElement>) => {
    const { onChange, signal, index } = this.props;
    const name = v.currentTarget.value;
    onChange({ ...signal, name }, index);
  };

  onLabelsChanged = (v: React.SyntheticEvent<HTMLInputElement>) => {
    const { onChange, signal, index } = this.props;
    const txt = v.currentTarget.value;
    const labels = txt ? parseLabels(txt) : undefined;
    onChange({ ...signal, labels }, index);
  };

  onUnitChanged = (v?: string) => {
    const { onChange, signal, index } = this.props;
    const copy = { ...signal };
    if (v) {
      if (!copy.config) {
        copy.config = {};
      }
      copy.config.unit = v;
    } else if (copy.config) {
      delete copy.config.unit;
    }
    onChange(copy, index);
  };

  onTrash = () => {
    const { onChange, index } = this.props;
    onChange(undefined, index);
  };

  renderLabels(signal: SignalField) {
    return (
      <InlineField label="Labels">
        <Input css="" defaultValue={formatLabels(signal.labels!)} onBlur={this.onLabelsChanged} placeholder="labels" />
      </InlineField>
    );
  }

  render() {
    const { signal, onAddExpr } = this.props;
    const exprs = [...standardWaves];

    let currentFn = exprs.find(e => e.value === signal.expr);
    if (!currentFn) {
      if (signal.expr) {
        currentFn = {
          value: signal.expr,
          label: signal.expr,
        };
      } else {
        currentFn = {};
      }
      exprs.push(currentFn);
    }

    return (
      <>
        <div className="gf-form">
          <InlineField label="Signal" labelWidth={10}>
            <Input width={10} css="" value={signal.name || ''} onChange={this.onNameChange} placeholder="Field" />
          </InlineField>
          <InlineField label="f(x)" grow={true}>
            <Select
              options={exprs}
              value={currentFn}
              onChange={this.onExprChange}
              placeholder="Function"
              allowCustomValue={true}
              isClearable={true}
              isSearchable={true}
              formatCreateLabel={txt => `fn: ${txt}`}
              onCreateOption={this.onNewExpr}
              menuPlacement="bottom"
            />
          </InlineField>
          {/* renderLabels(signal) */}
          <InlineField label="Unit">
            <UnitPicker value={signal.config?.unit} onChange={this.onUnitChanged} width={15} />
          </InlineField>
          {onAddExpr ? (
            <Button onClick={onAddExpr} variant="secondary">
              <Icon name="plus" />
            </Button>
          ) : (
            <Button onClick={this.onTrash} variant="secondary">
              <Icon name="trash-alt" />
            </Button>
          )}
        </div>
      </>
    );
  }
}
