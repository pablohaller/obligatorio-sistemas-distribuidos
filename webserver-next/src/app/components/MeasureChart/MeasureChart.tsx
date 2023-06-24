"use client";
import {
  CartesianGrid,
  Legend,
  Line,
  LineChart,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";

export interface MeasureChartData {
  pressure: number;
  datetime: string;
  sensor?: string;
  sector?: string;
}

interface Props {
  data: MeasureChartData;
}

const MeasureChart = ({ data }: Props) => {
  return (
    <ResponsiveContainer width="100%" height={200}>
      <LineChart
        data={data as any}
        margin={{
          top: 5,
          right: 30,
          left: 20,
          bottom: 5,
        }}
      >
        <Tooltip />
        <CartesianGrid strokeDasharray="3 3" />
        <XAxis dataKey="datetime" />
        <YAxis />
        <Legend />
        <Line
          name="Presión"
          type="monotone"
          dataKey="pressure"
          stroke="#0BA5E9"
          activeDot={{ r: 8 }}
        />
      </LineChart>
    </ResponsiveContainer>
  );
};

export default MeasureChart;
