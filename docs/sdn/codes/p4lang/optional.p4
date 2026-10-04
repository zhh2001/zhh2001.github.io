package Pipeline(/* 参数 */);
package Switch(Pipeline first, @optional Pipeline second);

Pipeline(/* 参数 */) ingress;
Switch(ingress) main;  // 一个只有单级流水线的交换机
